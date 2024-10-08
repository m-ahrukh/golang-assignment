package processor

type MatcherHandler struct {
	inputSrcA chan *Envelope
	inputSrcB chan *Envelope
	output    chan *Envelope

	application Matcher
}

type Matcher interface {
	Match(interface{}) JSONOutput
}

func NewMatcherHandler(inputSrcA, inputSrcB, output chan *Envelope, application Matcher) *MatcherHandler {
	return &MatcherHandler{
		inputSrcA:   inputSrcA,
		inputSrcB:   inputSrcB,
		output:      output,
		application: application,
	}
}

func (matcher *MatcherHandler) Handle() {
	for {
		select {
		case receivedA, ok := <-matcher.inputSrcA:
			if !ok {
				return
			}
			receivedA.Output = matcher.application.Match(receivedA.JsonInput)
			matcher.output <- receivedA
		case receivedB, ok := <-matcher.inputSrcB:
			if !ok {
				return
			}
			receivedB.Output = matcher.application.Match(receivedB.XmlInput)
			matcher.output <- receivedB
		}
	}
}
