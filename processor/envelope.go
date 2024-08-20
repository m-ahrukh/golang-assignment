package processor

type JSONInput struct {
	Id   string
	Kind string
}

type XMLInput struct {
}

type JSONOutput struct {
	Id   string
	Kind string
}

type Envelope struct {
	JsonInput JSONInput
	XmlInput  XMLInput

	Output JSONOutput
}
