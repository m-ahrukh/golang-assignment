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

	handler *MatcherHandler
}

func (mhf *MatcherHandlerFixture) Setup() {
	mhf.inputSrcA = make(chan *Envelope, 1024)
	mhf.inputSrcB = make(chan *Envelope, 1024)
	mhf.output = make(chan *Envelope, 1024)
	mhf.application = NewFakeMatcher()
	mhf.handler = NewMatcherHandler(mhf.inputSrcA, mhf.inputSrcB, mhf.output, mhf.application)
}

func (mhf *MatcherHandlerFixture) TestMatcherRecievesJSONInput() {
	envelope := mhf.enqueueJSONEnvelope("1", "Joined")
	close(mhf.inputSrcA)
	mhf.handler.Handle()

	mhf.AssertEqual("1", envelope.Output.Id)
	mhf.AssertEqual("Joined", envelope.Output.Kind)
	mhf.AssertEqual(envelope, <-mhf.output)
}

func (mhf *MatcherHandlerFixture) TestInputQueueDrained() {
	envelope1 := mhf.enqueueJSONEnvelope("2", "Joined")
	envelope2 := mhf.enqueueJSONEnvelope("3", "Joined")
	envelope3 := mhf.enqueueJSONEnvelope("4", "Joined")

	close(mhf.inputSrcA)
	mhf.handler.Handle()

	mhf.AssertEqual(envelope1, <-mhf.output)
	mhf.AssertEqual(envelope2, <-mhf.output)
	mhf.AssertEqual(envelope3, <-mhf.output)
}

func (mhf *MatcherHandlerFixture) enqueueJSONEnvelope(id string, kind string) *Envelope {
	envelope := &Envelope{
		JsonInput: JSONInput{
			Id:   id,
			Kind: kind,
		},
	}
	mhf.inputSrcA <- envelope
	return envelope
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

	output JSONOutput
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

	fakeMatcher.output = JSONOutput{Id: fakeMatcher.jsonInput.Id, Kind: fakeMatcher.jsonInput.Kind}
	return fakeMatcher.output
}
