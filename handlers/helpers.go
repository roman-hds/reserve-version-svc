package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/hashicorp/go-version"
)

//  PopulateCache retrieves existing build versions, sorts & picks the latest
func (app *Application) populateCache(branch string) error {
	buildVersions, err := app.ArtClient.GetBuilds(branch)
	if err != nil {
		app.ErrorLog.Println("retrieve current build directories %w\n", err.Error())
		return err
	}
	latestBuild := app.getLatest(buildVersions)
	app.InfoLog.Printf("latest build version for branch '%s' is '%s'", branch, latestBuild)

	app.InMemoryCache.Save(branch, latestBuild)
	return nil
}

// Returns the latest build version from provided slice
func (app *Application) getLatest(buildVersions []string) string {
	v1, _ := version.NewVersion("0.0.0.0")
	for _, buildNumber := range buildVersions {
		v2, err := version.NewVersion(buildNumber)
		if err != nil {
			app.ErrorLog.Println("error parsing version:", err.Error())
			continue
		}
		if v1.LessThan(v2) {
			v1 = v2
		}
	}
	return v1.String()
}

// Increment the build version by one
func (app *Application) incrementVersion(buildVersion string) string {
	var stringSlice []string

	v, _ := version.NewVersion(buildVersion)
	segments := v.Segments()
	segments[len(segments)-1]++
	for _, v := range segments {
		stringSlice = append(stringSlice, strconv.Itoa(v))
	}
	return strings.Join(stringSlice, ".")
}

// RespondWithError wraps the respondWithJSON function providing the handler’s ResponseWriter, an HTTP status code, and a payload to be marshaled
func (app *Application) RespondWithError(w http.ResponseWriter, code int, err error) {
	trace := fmt.Sprintf("%s\n%s\n", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	app.respondWithJSON(w, code, map[string]int{http.StatusText(code): code})
}

// respondWithJSON replies to the request with the formatted error message and HTTP code
func (app *Application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		app.ErrorLog.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
