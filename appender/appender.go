package appender

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
)

func FindEndOfType(fileName string, t string) (int, error) {
	pattern := fmt.Sprintf(`(?m)type (?<typeName>%s) struct[\w]?{[\s.\w]*(?<end>})`, t)
	reg := regexp.MustCompile(pattern)
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	index := reg.FindSubmatchIndex(content)
	endIndex := reg.SubexpIndex("end")
	insertPoint := index[endIndex*2] + 1
	return insertPoint, nil
}

func Append(fileName string, txt []byte, at int) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	before := content[:at]
	after := content[at:]
	buf := bytes.Buffer{}
	buf.Write(before)
	buf.WriteRune('\n')
	buf.WriteRune('\n')
	buf.Write(txt)
	buf.Write(after)
	return buf.Bytes(), nil
}
