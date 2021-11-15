package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"reserve-version/artifactory"
	"reserve-version/cache"
	"reserve-version/config"
	"reserve-version/handlers"
	"time"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	configFilePath := flag.String("cfg", "/data/app_configs/config-data.yaml", "Configuration YAML file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configFilePath)
	if err != nil {
		errorLog.Printf("failed to locate config file %s:", *configFilePath)
		return
	}

	artClient := &artifactory.ArtifactoryClient{
		Config: cfg,
		Logger: infoLog,
	}

	app := &handlers.Application{
		ArtClient:     artClient,
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		InMemoryCache: &cache.LatestBuildCache{Builds: make(map[string]string)},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/live", app.Healthz)
	mux.HandleFunc("/version", app.ReserveVersion)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      mux,
		ErrorLog:     app.ErrorLog,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	infoLog.Printf("Starting server on %s", srv.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
