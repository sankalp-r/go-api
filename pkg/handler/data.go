package handler

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	cache2 "github.com/sankalp-r/go-api/pkg/cache"
	"github.com/sankalp-r/go-api/pkg/model"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"
)

const (
	sortByView           = "views"
	sortByRelevanceScore = "relevanceScore"
	minLimit             = 1
	maxLimit             = 200
	sortQuery            = "sortKey"
	limitQuery           = "limit"
	invalidRequest       = "Invalid request"
	etagHeader           = "ETag"
	noneMatchHeader      = "If-None-Match"
)

// response for storing fetched data
type response struct {
	res *model.DataContainer
	err error
}

type DataHandler struct {
	cache      cache2.Storage
	httpClient *retryablehttp.Client
	seedUrl    []string
}

func NewDataHandler(seedUrl []string) *DataHandler {
	cache := cache2.NewStorage()
	httpClient := retryablehttp.NewClient()
	httpClient.RetryMax = 3
	httpClient.RetryWaitMin = 10 * time.Millisecond
	httpClient.RetryWaitMax = 50 * time.Millisecond
	httpClient.Logger = nil

	return &DataHandler{
		cache:      cache,
		httpClient: httpClient,
		seedUrl:    seedUrl,
	}
}

func (d *DataHandler) GetData(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("request received", zap.String("url", r.URL.String()))
	params := r.URL.Query()

	// sortQuery validation
	if params.Has(sortQuery) && !isSortKeyValid(params.Get(sortQuery)) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, invalidRequest)
		zap.L().Debug(invalidRequest, zap.Int("Status", http.StatusBadRequest))
		return
	}
	limit := 0
	//limitQuery validation
	if params.Has(limitQuery) {
		if ok, val := isLimitValid(params.Get(limitQuery)); ok {
			limit = val
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, invalidRequest)
			zap.L().Debug(invalidRequest, zap.Int("Status", http.StatusBadRequest))
			return
		}

	}

	responseChan := make(chan response, len(d.seedUrl))

	for _, url := range d.seedUrl {
		go func(url string) {
			data, err := d.fetchData(url)
			responseChan <- response{res: data, err: err}
		}(url)
	}

	var response []model.Data
	for i := 0; i < len(d.seedUrl); i++ {
		result := <-responseChan
		if result.err == nil {
			response = append(response, result.res.Data...)
		}
	}
	close(responseChan)

	resultContainer := sortData(response, params.Get(sortQuery), limit)
	bytes, err := json.Marshal(resultContainer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		zap.L().Error("internal server error", zap.Error(err))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(bytes)
	if err != nil {
		zap.L().Error("error in returning response", zap.Error(err))
	}

}

func GetSeedUrl() []string {
	return []string{
		"https://raw.githubusercontent.com/assignment132/assignment/main/duckduckgo.json",
		"https://raw.githubusercontent.com/assignment132/assignment/main/google.json",
		"https://raw.githubusercontent.com/assignment132/assignment/main/wikipedia.json",
	}
}

func (d *DataHandler) fetchData(url string) (*model.DataContainer, error) {
	request, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		zap.L().Error("error creating request", zap.Error(err))
		return nil, err
	}
	etag := d.cache.Get(url)
	if etag != nil {
		request.Header.Set(noneMatchHeader, etag.Key)
	}

	response, err := d.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	var body []byte

	switch response.StatusCode {
	case http.StatusOK:
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		d.cache.Set(url, cache2.Etag{Key: response.Header.Get(etagHeader), Data: body})
		zap.L().Debug("data fetch", zap.String("url", url), zap.Int("Status", http.StatusOK))
	case http.StatusNotModified:
		body = etag.Data
		zap.L().Debug("data fetch", zap.String("url", url), zap.Int("Status", http.StatusNotModified))
	default:
		zap.L().Error("failed to fetch data", zap.String("url", url), zap.Int("Status", response.StatusCode))
		return nil, fmt.Errorf("internal service error: %d", http.StatusInternalServerError)
	}

	defer response.Body.Close()
	var dataContainer model.DataContainer
	if err = json.Unmarshal(body, &dataContainer); err != nil {
		zap.L().Error("error in unmarshalling:", zap.Error(err))
		return nil, err
	}
	return &dataContainer, nil
}

func sortData(res []model.Data, sortType string, limit int) model.DataContainer {

	switch sortType {
	case sortByView:
		sort.Stable(model.DataByView(res))
	case sortByRelevanceScore:
		sort.Stable(model.DataByRelevanceScore(res))
	}

	if limit >= 1 && limit <= len(res) {
		res = res[:limit]
	}
	sortedResult := model.DataContainer{Data: res, Count: len(res)}
	return sortedResult
}

func isSortKeyValid(key string) bool {
	return key == sortByView || key == sortByRelevanceScore
}

func isLimitValid(key string) (bool, int) {
	limit, err := strconv.Atoi(key)
	if (err != nil) || !(limit >= minLimit && limit <= maxLimit) {
		return false, 0
	}
	return true, limit
}
