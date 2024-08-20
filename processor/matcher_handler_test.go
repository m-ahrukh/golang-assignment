package processor

import (
	"testing"

	"github.com/smarty/gunit"
)

func TestMatcherHandlerFixture(t *testing.T) {
	gunit.Run(new(MatcherHandlerFixture), t)
}

type MatcherHandlerFixture struct {
	*gunit.Fixture

	inputSrcA   chan *Envelope
	inputSrcB   chan *Envelope
	output      chan *Envelope
	application *FakeMatcher
	handler     *MatcherHandler
}

func (mhf *MatcherHandlerFixture) Setup() {
	mhf.inputSrcA = make(chan *Envelope, 10)
	mhf.inputSrcB = make(chan *Envelope, 10)
	mhf.output = make(chan *Envelope, 10)
	mhf.application = NewFakeMatcher()
	mhf.handler = NewMatcherHandler(mhf.inputSrcA, mhf.inputSrcB, mhf.output, mhf.application)
}

func (mhf *MatcherHandlerFixture) TestMatcherRecievesJSONInput() {
	envelope := &Envelope{
		JsonInput: JSONInput{
			Id:   "1",
			Kind: "Joined",
		},
	}

	mhf.application.output = JSONOutput{
		Id:   "1",
		Kind: "Joined",
	}
	mhf.inputSrcA <- envelope
	close(mhf.inputSrcA)
	mhf.handler.Handle()
	mhf.AssertEqual(envelope, <-mhf.output)
	mhf.AssertEqual("1", mhf.application.jsonInput.Id)
	mhf.AssertEqual("Joined", mhf.application.jsonInput.Kind)
}

func (mhf *MatcherHandlerFixture) TestMatcherRecievesXMLInput() {
	envelope := &Envelope{}

	mhf.inputSrcB <- envelope
	close(mhf.inputSrcB)
	mhf.handler.Handle()
	mhf.AssertEqual(envelope, <-mhf.output)
	mhf.AssertEqual(envelope.XmlInput, mhf.application.xmlInput)
}

// ////////////////////////////////////////////////////////////
type FakeMatcher struct {
	jsonInput JSONInput
	xmlInput  XMLInput
	output    JSONOutput
}

func NewFakeMatcher() *FakeMatcher {
	return &FakeMatcher{}
}

func (fakeMatcher *FakeMatcher) Match(value interface{}) JSONOutput {
	switch v := value.(type) {
	case JSONInput:
		fakeMatcher.jsonInput = v
	case XMLInput:
		fakeMatcher.xmlInput = v
	}

	return fakeMatcher.output
}
