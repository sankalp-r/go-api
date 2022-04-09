package router

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	for _, route := range routes {
		subRouter.HandleFunc(route.Pattern, route.HandlerFunc).Name(route.Name).Methods(route.Method)
	}

	return router
}
