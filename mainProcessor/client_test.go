package processor

import (
	"fmt"
	"testing"
)

func SKipTestJSON(t *testing.T) {
	data, err := FetchData("JSON")

	if err != nil {
		t.Errorf("Got an error: %s", err)
	}

	fmt.Println("JSON Id:", data)
}

func SkipTestXML(t *testing.T) {
	data, err := FetchData("XML")

	if err != nil {
		t.Errorf("Got an error: %s", err)
	}

	fmt.Println("XML Id:", data)
}

func TestMatch(t *testing.T) {
	// jsonId, _ := FetchData("JSON")
	// xmlId, _ := FetchData("XML")

	// if Matcher(jsonId, xmlId) {
	// 	fmt.Println("Joined")
	// }
	// println("Orphaned")

	FetchAndMatcher()
}
