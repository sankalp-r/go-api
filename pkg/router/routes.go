package router

import (
	"github.com/sankalp-r/go-api/pkg/handler"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []Route{
	{
		"GetData",
		"GET",
		"/data",
		handler.NewDataHandler(handler.GetSeedUrl()).GetData,
	},
}
