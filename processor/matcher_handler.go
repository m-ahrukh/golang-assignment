package processor

type MatcherHandler struct {
	inputSrcA chan interface{}
	inputSrcB chan interface{}
	output    chan interface{}

	application Matcher
}

type Matcher interface {
	Match(interface{})
}

func NewMatcherHandler(inputSrcA, inputSrcB, output chan interface{}, application Matcher) *MatcherHandler {
	return &MatcherHandler{
		inputSrcA:   inputSrcA,
		inputSrcB:   inputSrcB,
		output:      output,
		application: application,
	}
}

func (matcher *MatcherHandler) Handle() {
	// recieved := <-matcher.inputSrcA
	// matcher.application.Match(recieved)
	// matcher.output <- recieved

	for {
		select {
		case receivedA, ok := <-matcher.inputSrcA:
			if !ok {
				return
			}
			matcher.application.Match(receivedA)
			matcher.output <- receivedA
		case receivedB, ok := <-matcher.inputSrcB:
			if !ok {
				return
			}
			matcher.application.Match(receivedB)
			matcher.output <- receivedB
		}
	}
}
