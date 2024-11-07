package appender

import "bytes"

type Appender struct {
	beforeText   *bytes.Buffer
	afterText    *bytes.Buffer
	FileName     string
	TextToAppend string
	AfterType    string
}

func NewAppender() *Appender {
	a := &Appender{}
	a.beforeText = &bytes.Buffer{}
	a.afterText = &bytes.Buffer{}
	return a
}
