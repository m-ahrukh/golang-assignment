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

	inputSrcA   chan interface{}
	inputSrcB   chan interface{}
	output      chan interface{}
	application *FakeMatcher
	handler     *MatcherHandler
}

func (mhf *MatcherHandlerFixture) Setup() {
	mhf.inputSrcA = make(chan interface{}, 10)
	mhf.inputSrcB = make(chan interface{}, 10)
	mhf.output = make(chan interface{}, 10)
	mhf.application = NewFakeMatcher()
	mhf.handler = NewMatcherHandler(mhf.inputSrcA, mhf.inputSrcB, mhf.output, mhf.application)
}

func (mhf *MatcherHandlerFixture) TestMatcherRecievesInputA() {
	mhf.inputSrcA <- 1
	close(mhf.inputSrcA)
	mhf.handler.Handle()
	mhf.AssertEqual(1, <-mhf.output)
	mhf.AssertEqual(1, mhf.application.inputSrc)

	mhf.inputSrcB <- 2
	close(mhf.inputSrcB)
	mhf.handler.Handle()
	mhf.AssertEqual(2, <-mhf.output)
	mhf.AssertEqual(2, mhf.application.inputSrc)
}

// ////////////////////////////////////////////////////////////
type FakeMatcher struct {
	inputSrc interface{}
}

func NewFakeMatcher() *FakeMatcher {
	return &FakeMatcher{}
}

func (fakeMatcher *FakeMatcher) Match(value interface{}) {
	fakeMatcher.inputSrc = value
}
