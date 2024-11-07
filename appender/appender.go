package appender

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
)

func Append(fileName string, txt []byte, afterType string) ([]byte, error) {
	pattern := fmt.Sprintf(`(?m)type (?<typeName>%s) struct[\w]?{[\s.\w]*(?<end>})`, afterType)
	reg := regexp.MustCompile(pattern)
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	index := reg.FindSubmatchIndex(content)
	endIndex := reg.SubexpIndex("end")
	insertPoint := index[endIndex*2] + 1
	before := content[:insertPoint]
	after := content[insertPoint:]
	buf := bytes.Buffer{}
	buf.Write(before)
	buf.WriteRune('\n')
	buf.WriteRune('\n')
	buf.Write(txt)
	buf.Write(after)
	return buf.Bytes(), nil
}
