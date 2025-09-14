package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/samirhembrom/httpfromtcp/internal/headers"
)

type StatusCode int

type WriteState int

type Writer struct {
	w          io.Writer
	writeState WriteState
}

const (
	WritingStatusLine WriteState = iota
	WritingHeaders
	WritingBody
)

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

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writeState != WritingStatusLine {
		return fmt.Errorf("writes not in correct order")
	}

	err := WriteStatusLine(w.w, statusCode)
	if err != nil {
		return err
	}
	w.writeState = WritingHeaders
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writeState != WritingHeaders {
		return fmt.Errorf("writes not in correct order")
	}
	err := WriteHeaders(w.w, headers)
	if err != nil {
		return err
	}
	w.writeState = WritingBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writeState != WritingBody {
		return 0, fmt.Errorf("writes not in correct order")
	}
	n, err := w.w.Write(p)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:          w,
		writeState: WritingStatusLine,
	}
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	bytesLenInt := len(p)
	bytesLenHex := fmt.Sprintf("%x", bytesLenInt)
	data := append([]byte(bytesLenHex+"\r\n"), p...)
	data = append(data, []byte("\r\n")...)
	_, err := w.w.Write(data)
	if err != nil {
		return 0, err
	}
	return bytesLenInt, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writeState != WritingBody {
		return 0, fmt.Errorf("writes not in correct order")
	}
	body := "0\r\n\r\n"
	n, err := w.w.Write([]byte(body))
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	w.w.Write([]byte("0\r\n"))
	if v, ok := h.Get("x-content-sha256"); ok {
		if _, err := w.w.Write([]byte("X-Content-Sha256: " + v + "\r\n")); err != nil {
			return err
		}
	}
	if v, ok := h.Get("x-content-length"); ok {
		if _, err := w.w.Write([]byte("X-Content-Length: " + v + "\r\n")); err != nil {
			return err
		}
	}
	_, err := w.w.Write([]byte("\r\n"))
	return err
}
