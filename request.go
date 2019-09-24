package auth

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	// HTTPClient is an overrideable HTTP client for sending requests
	HTTPClient = &http.Client{Timeout: 5 * time.Second}

	// PlanningCenterScheme is the overrideable HTTP scheme for URLs generated
	PlanningCenterScheme = "https"

	// PlanningCenterHost is the overrideable HTTP host for URLs generated
	PlanningCenterHost = "api.planningcenteronline.com"
)

func sendRequest(req *http.Request, logger *log.Logger) ([]byte, error) {
	logger.Printf("[%s] %s", req.Method, req.URL)

	response, err := HTTPClient.Do(req)
	if err != nil {
		return []byte{}, err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return []byte{}, fmt.Errorf("HTTP Error: %s", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	return body, nil
}

func apiURL(path string) *url.URL {
	return &url.URL{
		Scheme: PlanningCenterScheme,
		Host:   PlanningCenterHost,
		Path:   path,
	}
}

func redirectURI(port int) *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%d", port),
		Path:   "/callback",
	}
}
