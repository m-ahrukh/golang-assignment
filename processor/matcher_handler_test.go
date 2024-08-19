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

func (mhf *MatcherHandlerFixture) TestMatcherRecievesInput() {
	envelope := &Envelope{}
	mhf.inputSrcA <- envelope
	close(mhf.inputSrcA)
	mhf.handler.Handle()
	mhf.AssertEqual(envelope, <-mhf.output)
	mhf.AssertEqual(envelope.JsonInput, mhf.application.jsonInput)
	// mhf.AssertEqual(1, <-mhf.output)
	// mhf.AssertEqual(1, mhf.application.inputSrc)

	mhf.inputSrcB <- envelope
	close(mhf.inputSrcB)
	mhf.handler.Handle()
	mhf.AssertEqual(envelope, <-mhf.output)
	mhf.AssertEqual(envelope.XmlInput, mhf.application.xmlInput)
	// mhf.AssertEqual(2, <-mhf.output)
	// mhf.AssertEqual(2, mhf.application.inputSrc)
}

// ////////////////////////////////////////////////////////////
type FakeMatcher struct {
	jsonInput JSONInput
	xmlInput  XMLInput
}

func NewFakeMatcher() *FakeMatcher {
	return &FakeMatcher{}
}

func (fakeMatcher *FakeMatcher) Match(value interface{}) {
	switch v := value.(type) {
	case JSONInput:
		fakeMatcher.jsonInput = v
	case XMLInput:
		fakeMatcher.xmlInput = v
	}
}
