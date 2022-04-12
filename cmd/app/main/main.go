package main

import (
	"context"
	"flag"
	muxrouter "github.com/sankalp-r/go-api/pkg/router"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	var address string
	var enableDebugLog bool
	flag.StringVar(&address, "address", "8080", "HTTP Server Address")
	flag.BoolVar(&enableDebugLog, "enableDebug", false, "Enable debug level log")
	flag.Parse()
	config := zap.NewProductionConfig()
	if enableDebugLog {
		config.Level.SetLevel(zap.DebugLevel)
	}
	logger, _ := config.Build()
	zap.ReplaceGlobals(logger)
	router := muxrouter.NewRouter()
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	srv := &http.Server{
		Addr:    ":" + address,
		Handler: router,
	}

	go func() {
		zap.L().Info("Server started..", zap.String("port", address))
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-stop

	zap.L().Info("Shutting down server..")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

}
