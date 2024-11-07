package appender

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"k8s.io/gengo/types"
)

type Appender struct {
	IfaceMethods      map[string]*types.Type
	TargetTypeMethods [][]string
}

func New(methods map[string]*types.Type, targetMeth [][]string) *Appender {
	a := &Appender{}
	a.IfaceMethods = methods
	a.TargetTypeMethods = targetMeth
	return a
}

func (a *Appender) FindEndOfType(fileName string, t string) (int, error) {
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

func (a *Appender) Append(fileName string, at int, typeName string, ifaceName string) error {
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
	txt := bytes.Buffer{}
outer:
	for n, v := range a.IfaceMethods {
		for _, brn := range a.TargetTypeMethods {
			if strings.ReplaceAll(typeName, "*", "") == strings.ReplaceAll(brn[1], "*", "") && brn[2] == n {
				txt.WriteString(brn[0])
				txt.WriteRune('\n')
				continue outer
			}
		}
		toRunes := []rune(strings.ToLower(typeName))
		recver := fmt.Sprintf("func(%s %s)%s", string(toRunes[0]), typeName, n)
		args := strings.Builder{}
		args.WriteRune('(')
		for ii := 0; ii < len(v.Signature.ParameterNames); ii++ {
			args.WriteString(v.Signature.ParameterNames[ii])
			args.WriteString(" ")
			args.WriteString(v.Signature.Parameters[ii].String())
			if ii != len(v.Signature.ParameterNames)-1 {
				args.WriteRune(',')
			}
		}
		args.WriteString(")")
		result := strings.Builder{}
		resLen := len(v.Signature.Results)
		if resLen == 1 {
			result.WriteString(v.Signature.Results[0].Name.String())
		} else if resLen > 1 {
			result.WriteRune('(')
			for ii := 0; ii < resLen; ii++ {
				result.WriteString(v.Signature.Results[ii].String())
				if ii != resLen-1 {
					result.WriteRune(',')
				}
			}
			result.WriteRune(')')
		}

		txt.WriteString(fmt.Sprintf("%s%s%s{\n  panic(\"to be implemented\")\n}\n", recver, args.String(), result.String()))
	}
	buf.Write(txt.Bytes())
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

func (a Appender) DeleteLastAppend(fileName string, typeName string, ifaceName string) error {
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
	}
	return nil
}
