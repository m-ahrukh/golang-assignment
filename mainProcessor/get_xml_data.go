package processor

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type MSG struct {
	ID   *ID    `xml:"id"`
	Done string `xml:"done"`
}

type ID struct {
	Value string `xml:"value,attr"`
}

func XmlUrl() string {
	endpoint := "/source/b"
	return MakeUrl(baseUrl, endpoint)
}

func getXMLRequest(url string) (*http.Request, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	return request, nil
}

func setXMLHeaders(request http.Request) {
	request.Header.Set("Content-Type", "application/xml")
}

func getXMLResponse(request *http.Request) (*http.Response, error) {
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	return response, nil
}

func parsingXMLResponseIntoBytes(response *http.Response) ([]byte, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, nil
}

func parseXMLData(body []byte) (*MSG, error) {
	var msg MSG
	err := xml.Unmarshal(body, &msg)
	if err != nil {
		if strings.Contains(string(body), "</foo>") {
			return nil, errors.New("malformed XML detected")
		}
		return nil, fmt.Errorf("failed to unmarshal XML: %w", err)
	}
	return &msg, nil
}

func FetchXMLData(client *http.Client, url string) (*MSG, error) {
	request, err := getXMLRequest(url)
	if err != nil {
		return nil, err
	}

	setXMLHeaders(*request)

	response, err := getXMLResponse(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	body, err := parsingXMLResponseIntoBytes(response)
	if err != nil {
		return nil, err
	}

	data, err := parseXMLData(body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func HandleXMLResponse(msg *MSG) (*MSG, error) {
	if msg.Done != "" {
		fmt.Println("processing completed.")
		return msg, nil
	}
	if msg.ID != nil {
		return msg, nil
	}
	return nil, errors.New("unknown response format")
}
