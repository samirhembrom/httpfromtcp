package response

import (
	"io"
	"strconv"

	"github.com/samirhembrom/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK            StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case 200:
		w.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case 400:
		w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case 500:
		w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	default:
		codeStr := strconv.Itoa(int(statusCode))
		statusLine := "HTTP/1.1 " + codeStr + " \r\n"
		w.Write([]byte(statusLine))
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	str := strconv.Itoa(contentLen)
	headers.Set("Content-Length", str)
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		header := k + ": " + v + "\r\n"
		_, err := w.Write([]byte(header))
		if err != nil {
			return err
		}
	}
	w.Write([]byte("\r\n"))
	return nil
}
