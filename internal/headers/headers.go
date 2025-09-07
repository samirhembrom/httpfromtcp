package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return idx + 2, true, nil
	}
	fieldLineText := string(data[:idx])
	fieldLine, err := fieldLikeFromString(fieldLineText)
	if err != nil {
		return 0, false, err
	}
	sepIdx := bytes.Index([]byte(fieldLine), []byte(":"))
	h[fieldLine[:sepIdx]] = fieldLine[sepIdx+1:]
	return idx + 2, false, nil
}

func fieldLikeFromString(fieldLineText string) (string, error) {
	fieldLineText = strings.TrimSpace(fieldLineText)
	idx := bytes.Index([]byte(fieldLineText), []byte(":"))
	if idx == -1 {
		return "", errors.New("Invalid header")
	}
	fieldName := fieldLineText[:idx]
	if strings.Contains(fieldName, " ") {
		return "", errors.New("Field name contains space")
	}
	fieldValue := fieldLineText[idx+1:]
	fieldValue = strings.TrimSpace(fieldValue)
	return fieldName + ":" + fieldValue, nil
}
