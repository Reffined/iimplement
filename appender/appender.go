package appender

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
)

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

func (a *Appender) Append(fileName string, afterType string) error {
	pattern := `(?m)type (?<typeName>.+?) struct[\w]?{[\s.\w]*(?<end>})`
	reg := regexp.MustCompile(pattern)
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	index := reg.FindSubmatchIndex(content)
	for i := 2; i < len(index); i += 2 {
		l := index[i+1] - index[i]
		reader := io.NewSectionReader(file, int64(index[i]), int64(l))
		buf := make([]byte, l)
		_, err := reader.Read(buf)
		if err != nil {
			return err
		}
		fmt.Println(string(buf))

	}

	return nil
}
