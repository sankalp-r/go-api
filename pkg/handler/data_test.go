package handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testData = `{"data":[{"url":"www.test.com/abc1","views":1000,"relevanceScore":0.1},{"url":"www.test.com/abc2","views":3000,"relevanceScore":0.3},{"url":"www.test.com/abc3","views":2000,"relevanceScore":0.2}]}`

func TestGetData(t *testing.T) {
	testSeedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, testData)
	}))
	defer testSeedServer.Close()
	testSeedUrl := testSeedServer.URL

	testCases := []struct {
		name                 string
		request              string
		expectedResponse     string
		expectedResponseCode int
	}{
		{
			name:                 "without sort",
			request:              "/api/v1/data",
			expectedResponse:     testData,
			expectedResponseCode: 200,
		},
		{
			name:                 "with sort by views",
			request:              "/api/v1/data?sortKey=views",
			expectedResponse:     `{"data":[{"url":"www.test.com/abc1","views":1000,"relevanceScore":0.1},{"url":"www.test.com/abc3","views":2000,"relevanceScore":0.2},{"url":"www.test.com/abc2","views":3000,"relevanceScore":0.3}]}`,
			expectedResponseCode: 200,
		},
		{
			name:                 "with sort by relevance score",
			request:              "/api/v1/data?sortKey=relevanceScore",
			expectedResponse:     `{"data":[{"url":"www.test.com/abc1","views":1000,"relevanceScore":0.1},{"url":"www.test.com/abc3","views":2000,"relevanceScore":0.2},{"url":"www.test.com/abc2","views":3000,"relevanceScore":0.3}]}`,
			expectedResponseCode: 200,
		},
		{
			name:                 "with sort by relevance score and limit",
			request:              "/api/v1/data?sortKey=relevanceScore&limit=2",
			expectedResponse:     `{"data":[{"url":"www.test.com/abc1","views":1000,"relevanceScore":0.1},{"url":"www.test.com/abc3","views":2000,"relevanceScore":0.2}]}`,
			expectedResponseCode: 200,
		},
		{
			name:                 "Request with invalid sortkey",
			request:              "/api/v1/data?sortKey=rv",
			expectedResponse:     "Invalid request",
			expectedResponseCode: 400,
		},
	}

	testDataHandler := NewDataHandler([]string{testSeedUrl})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", tc.request, nil)
			w := httptest.NewRecorder()
			testDataHandler.GetData(w, r)
			response := w.Result()
			body, err := ioutil.ReadAll(response.Body)
			assert.Equal(t, err, nil)
			assert.Equal(t, body, []byte(tc.expectedResponse))
			fmt.Println(string(body))
			assert.Equal(t, response.StatusCode, tc.expectedResponseCode)
		})

	}

}
