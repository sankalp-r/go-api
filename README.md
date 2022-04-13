# go-api

* To download and sync the dependencies run the command `make deps-update` from `Makefile`.
* To run the API locally, run the command `make run`.
* To run unit-tests, run `make test`
* To deploy the API on Kubernetes, run command `make deploy`
* To build and publish the docker image, run `make image` and then `make push`(update the REGISTRY and image repo/tag accordingly)


## Design details
* For retry mechanism, `retryableHTTPclient` has been used.([go-retryablehttp](github.com/hashicorp/go-retryablehttp)).
It provides retry with backoff features.
* Using the ETAGs from query-URL response headers , query-URL response is cached in-memory for better performance.
* If some query-URL is returning error, then it will be excluded from final response.
* If the API request does not have `sortKey` parameter then unsorted response will be returned. 
* If the API request does not have `limit` parameter then all the response will be returned.
* Port for running the service can be changed using command line arguments:
`go run main.go --address=8081`
* Enable debug log-level by specifying `--enableDebug` argument while running the service.


## Requirement checklist
- [x] The server should query URLs concurrently.
- [x] Server should have re-try mechanism and error handling on failures while querying URLs.
- [x] Code should have unit tests
- [x] Included manifest files for deployment and Dockerfile.


## Enhancement scope
* Swagger docs can be added for better documentation of API.
* To control the concurrency-limit, worker-pool can be used.
* API can be made secure, to make it production-grade.