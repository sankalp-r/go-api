package main

import (
	"flag"
	muxrouter "github.com/sankalp-r/go-api/pkg/router"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	var address string
	var enableDebugLog bool
	flag.StringVar(&address, "address", ":8080", "HTTP Server Address")
	flag.BoolVar(&enableDebugLog, "enableDebug", false, "Enable debug level log")
	config := zap.NewProductionConfig()
	if enableDebugLog {
		config.Level.SetLevel(zap.DebugLevel)
	}
	logger, _ := config.Build()
	zap.ReplaceGlobals(logger)
	router := muxrouter.NewRouter()
	log.Fatal(http.ListenAndServe(address, router))
}
