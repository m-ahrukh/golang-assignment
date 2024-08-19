package processor

type JSONInput struct {
}

type XMLInput struct {
}

type Envelope struct {
	JsonInput JSONInput
	XmlInput  XMLInput
}
