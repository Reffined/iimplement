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

func Append(fileName string, txt []byte, at int, typeName string, ifaceName string) error {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	before := content[:at]
	after := content[at:]
	buf := bytes.Buffer{}
	buf.Write(before)
	buf.WriteRune('\n')
	buf.WriteRune('\n')
	buf.WriteString(fmt.Sprintf("// +iipml:%s:%s:begin\n", typeName, ifaceName))
	buf.Write(txt)
	buf.WriteString(fmt.Sprintf("// +iipml:%s:%s:end\n", typeName, ifaceName))
	buf.Write(after)
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = file.Write(buf.Bytes())
	if err != nil {
		return err
	}
	file.Truncate(int64(buf.Len()))
	return nil
}

func DeleteLastAppend(fileName string, typeName string, ifaceName string) error {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	begin := regexp.MustCompile(fmt.Sprintf(`// \+iipml:%s:%s:begin\n`, typeName, ifaceName))
	end := regexp.MustCompile(fmt.Sprintf(`// \+iipml:%s:%s:end\n`, typeName, ifaceName))
	beginIndex := begin.FindIndex(content)
	endIndex := end.FindIndex(content)
	if beginIndex != nil && endIndex != nil {
		buf := bytes.Buffer{}
		buf.Write(content[:beginIndex[0]])
		buf.Write(content[endIndex[1]:])
		_, err := file.Seek(0, 0)
		if err != nil {
			return err
		}
		_, err = file.Write(buf.Bytes())
		if err != nil {
			return err
		}
		file.Truncate(int64(buf.Len()))
		fmt.Println(buf.String())
	}
	return nil
}
