package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var client = http.Client{}

type ResultData struct {
	Kind string `json:"kind"`
	Id   string `json:"id"`
}

const baseUrl = "http://localhost:7299"

func MakeUrl(baseUrl, endpoint string) string {
	return fmt.Sprintf("%s%s", baseUrl, endpoint)
}

func PostResult(resultType, id string) error {
	url := MakeUrl(baseUrl, "/sink/a")

	result := ResultData{
		Kind: resultType,
		Id:   id,
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create POST request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send POST request:  %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func FetchData(responseType string) (string, error) {
	switch responseType {
	case "JSON":
		url := JsonURL()
		data, err := FetchJSONData(&client, url)
		if err != nil {
			if err.Error() == "received done response" {
				fmt.Println("Done message received in JSON")
				return "done", nil
			} else if err.Error() == "malformed json detected" {
				fmt.Println("Malformed JSON received.")
			} else {
				fmt.Printf("Error fetching JSON data: %v\n", err)
			}
			return "", err
		}
		return data.Id, nil
	case "XML":
		url := XmlUrl()
		data, err := FetchXMLData(&client, url)
		if err != nil {
			if err.Error() == "malformed XML detected" {
				fmt.Println("Malformed XML received.")
			} else {
				fmt.Printf("Error fetching XML data: %v\n", err)
			}
			return "", err
		} else {
			parsedMsg, err := HandleXMLResponse(data)
			if err != nil {
				return "", err
			} else {
				if parsedMsg.ID.Value != "" {
					return parsedMsg.ID.Value, nil
				} else {
					return "<done/>", nil
				}
			}
		}
	}
	return "", fmt.Errorf("invalid response type")
}

func FetchAndMatcher() {
	var wg sync.WaitGroup
	jsonChan := make(chan string, 2048)
	xmlChan := make(chan string, 2048)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			jsonId, err := FetchData("JSON")
			if err != nil {
				jsonChan <- ""
				continue
			}
			if jsonId == "done" {
				break
			}
			fmt.Println("JSON Id: ", jsonId)
			jsonChan <- jsonId
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			xmlId, err := FetchData("XML")
			if err != nil {
				xmlChan <- ""
				continue
			}

			if xmlId == "<done/>" {
				fmt.Println("Received XML done msg")
				break
			}

			fmt.Println("XML Id: ", xmlId)
			xmlChan <- xmlId

		}
	}()

	go func() {
		wg.Wait()
		close(jsonChan)
		close(xmlChan)
	}()

	jsonResults := make(map[string]bool)
	xmlResults := make(map[string]bool)

	for id := range jsonChan {
		if id != "" {
			jsonResults[id] = true
		}
	}

	for id := range xmlChan {
		if id != "" {
			xmlResults[id] = true
		}
	}

	fmt.Println("JSON Result: ", jsonResults)
	fmt.Println("XML Result: ", xmlResults)

	for id := range jsonResults {
		if xmlResults[id] {
			fmt.Printf("joined %s\n", id)
			err := PostResult("Joined", id)
			if err != nil {
				fmt.Printf("Error posting joined result: %v\n", err)
			}
		} else {
			fmt.Printf("orphaned %s\n", id)
			err := PostResult("Orphaned", id)
			if err != nil {
				fmt.Printf("Error posting oephaned result: %v\n", err)
			}
		}
	}

	for id := range xmlResults {
		if !jsonResults[id] {
			fmt.Printf("orphaned %s\n", id)
			err := PostResult("Orphaned", id)
			if err != nil {
				fmt.Printf("Error posting orphaned result:%v\n", err)
			}
		}
	}
}
