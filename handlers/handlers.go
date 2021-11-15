package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reserve-version/artifactory"
	"reserve-version/cache"
)

type Application struct {
	ArtClient     *artifactory.ArtifactoryClient
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InMemoryCache *cache.LatestBuildCache
}

// healthz is a liveness probe.
func (app *Application) Healthz(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

// reserveVersion reads value from the in-memory cache, increments it, creates a new directory with the new build version and saves updated version back to in-memory cache
func (app *Application) ReserveVersion(w http.ResponseWriter, req *http.Request) {
	branch := req.URL.Query().Get("branch")
	if len(branch) == 0 {
		err := errors.New("branch is not specified")
		app.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	// if specified branch doesn't have associated buildVersion -> populate it first
	if !app.InMemoryCache.HasKey(branch) {
		app.InfoLog.Println("in-memory cache created a key-value pair for provided branch")
		app.populateCache(branch)
	}

	latestBuild := app.InMemoryCache.Read(branch)
	// increment by 1
	newBuildVersion := app.incrementVersion(latestBuild)

	err := app.ArtClient.CreateBuildDir(branch, newBuildVersion)
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	app.InfoLog.Println("in-memory cache saved a key-value pair for provided branch")
	app.InMemoryCache.Save(branch, newBuildVersion)
	fmt.Fprintf(w, "%s", newBuildVersion)
}
