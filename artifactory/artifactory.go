package artifactory

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reserve-version/config"
	"strings"
	"time"
)

type ArtifactoryClient struct {
	Config *config.AppConfig
	Logger *log.Logger
}

type Build struct {
	Uri string `json:"uri"`
}

type BuildVersions struct {
	Builds []Build `json:"children"`
}

// GetBuilds returns a slice of existing build versions
func (art *ArtifactoryClient) GetBuilds(branch string) ([]string, error) {
	var buildVersions BuildVersions
	requestURL := art.Config.CurrentBuilds + branch + "/"

	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("get request: %w", err)
	}

	byteResponse, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	err = json.Unmarshal(byteResponse, &buildVersions)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response from Artifactory %s: %w", byteResponse, err)
	}

	if len(buildVersions.Builds) == 0 {
		return nil, errors.New("build directory is empty")
	}

	var buildsSlice []string
	for _, child := range buildVersions.Builds {
		buildsSlice = append(buildsSlice, strings.TrimPrefix(child.Uri, "/"))
	}
	return buildsSlice, nil
}

// Creates new build directory
func (art *ArtifactoryClient) CreateBuildDir(branch, buildVersion string) error {
	var responseStruct interface{} // for logging response
	targetURL := art.Config.CreateDir + branch + "/" + buildVersion + "/"

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodPut, targetURL, nil)
	if err != nil {
		return fmt.Errorf("error %s", err.Error())
	}
	req.SetBasicAuth(art.Config.User, art.Config.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot create build directory %s %w", buildVersion, err)
	}

	if err = unmarshalJSON(resp, &responseStruct); err != nil {
		return err
	}
	art.Logger.Printf("{'create dir request': ['branch:%s', 'buildVersion:%s']}\n Response: %s\n", branch, buildVersion, responseStruct)
	return nil
}

// Helper method for unmarshaling response
func unmarshalJSON(resp *http.Response, responseStruct *interface{}) error {
	byteResponse, _ := ioutil.ReadAll(resp.Body)
	// The error returned from response.Body.Close() can be ignored as the byte buffer of the inbound request body is wrapped in a nopCloser,
	// which has a Close method that always returns nil
	defer resp.Body.Close()

	// Unmarshal
	err := json.Unmarshal(byteResponse, &responseStruct)
	if err != nil {
		return fmt.Errorf("error unmarshaling response from Artifactory %w", err)
	}
	return nil
}
