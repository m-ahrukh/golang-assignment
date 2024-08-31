package processor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type JSONData struct {
	Status string `json:"status"`
	Id     string `json:"id"`
}

func JsonURL() string {
	endpoint := "/source/a"
	return MakeUrl(baseUrl, endpoint)
}

func getJSONRequest(url string) (*http.Request, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	return request, nil
}

func setJsonHeaders(request http.Request) {
	request.Header.Set("Content-Type", "application/json")
}

func getJSONResponse(request *http.Request) (*http.Response, error) {
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	return response, nil
}

func parsingJSONResponseIntoBytes(response *http.Response) ([]byte, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, nil
}

func parseJSONData(body []byte) (*JSONData, error) {
	var jsonData JSONData

	err := json.Unmarshal(body, &jsonData)
	if err != nil {
		if strings.Contains(string(body), "[") {
			return nil, errors.New("malformed json detected")
		}
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	if jsonData.Status == "done" {
		return nil, fmt.Errorf("received done response")
	}
	return &jsonData, nil
}

func FetchJSONData(client *http.Client, url string) (*JSONData, error) {
	request, err := getJSONRequest(url)
	if err != nil {
		return nil, err
	}

	setJsonHeaders(*request)

	response, err := getJSONResponse(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	body, err := parsingJSONResponseIntoBytes(response)
	if err != nil {
		return nil, err
	}
	jsonData, err := parseJSONData(body)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
