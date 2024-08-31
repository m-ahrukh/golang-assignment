package processor

import (
	"net/http"
	"testing"

	"github.com/smarty/gunit"
)

func TestMatcherFixture(t *testing.T) {
	gunit.Run(new(MatcherFixture), t)
}

type MatcherFixture struct {
	*gunit.Fixture

	client  *FakeHTTPClient
	matcher *WebServerMatcher
}

func (mf *MatcherFixture) Setup() {
	mf.client = &FakeHTTPClient{}
	mf.matcher = NewWebServerMatcher(mf.client)
}

func NewWebServerMatcher(client HTTPClient) *WebServerMatcher {
	return &WebServerMatcher{
		client: client,
	}
}

func (mf *MatcherFixture) TestRequestComosedProperly() {
	jsonInput := JSONInput{
		Id:   "1",
		Kind: "Joined",
	}

	xmlInput := XMLInput{}
	mf.matcher.Match(jsonInput, xmlInput)

	mf.AssertEqual("GET", mf.client.request.Method)

}

// /////////////////////////////////////////////////////////////
type FakeHTTPClient struct {
	request *http.Request
}

func (client *FakeHTTPClient) Do(request *http.Request) (*http.Response, error) {
	client.request = request
	return nil, nil
}
