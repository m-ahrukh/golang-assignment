package processor

import (
	"net/http"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}
type WebServerMatcher struct {
	client HTTPClient
}

func (server *WebServerMatcher) Match(JSONInput, XMLInput) JSONOutput {
	request, _ := http.NewRequest("GET", "", nil)
	server.client.Do(request)
	return JSONOutput{}
}
