package handler

import (
	"encoding/json"
	"fmt"
	cache2 "github.com/sankalp-r/go-api/pkg/cache"
	"github.com/sankalp-r/go-api/pkg/model"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
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

type response struct {
	res *model.DataContainer
	err error
}

var cache cache2.Storage = cache2.NewStorage()

func GetData(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	if params.Has(sortQuery) && !isSortKeyValid(params.Get(sortQuery)) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, invalidRequest)
		zap.L().Debug(invalidRequest, zap.Int("Status", http.StatusBadRequest))
		return
	}
	limit := 0
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

	seedUrl := getSeedUrl()
	responseChan := make(chan response, len(seedUrl))

	for _, url := range seedUrl {
		go func(url string) {
			data, err := fetchData(url)
			responseChan <- response{res: data, err: err}
		}(url)
	}

	var response []model.Data
	for i := 0; i < len(seedUrl); i++ {
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
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)

}

func getSeedUrl() []string {
	return []string{
		"https://raw.githubusercontent.com/assignment132/assignment/main/duckduckgo.json",
		"https://raw.githubusercontent.com/assignment132/assignment/main/google.json",
		"https://raw.githubusercontent.com/assignment132/assignment/main/wikipedia.json",
	}
}

func fetchData(url string) (*model.DataContainer, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	etag := cache.Get(url)
	if etag != nil {
		request.Header.Set(noneMatchHeader, etag.Key)
	}

	response, err := http.DefaultClient.Do(request)
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
		cache.Set(url, cache2.Etag{Key: response.Header.Get(etagHeader), Data: body})
		zap.L().Debug("data fetch", zap.String("url", url), zap.Int("Status", http.StatusOK))
	case http.StatusNotModified:
		body = etag.Data
		zap.L().Debug("data fetch", zap.String("url", url), zap.Int("Status", http.StatusNotModified))
	default:
		zap.L().Error("failed to fetch data", zap.String("url", url), zap.Int("Status", response.StatusCode))
	}

	defer response.Body.Close()
	var dataContainer model.DataContainer
	if err = json.Unmarshal(body, &dataContainer); err != nil {
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
	sortedResult := model.DataContainer{Data: res}
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
