package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/samirhembrom/httpfromtcp/internal/headers"
	"github.com/samirhembrom/httpfromtcp/internal/request"
	"github.com/samirhembrom/httpfromtcp/internal/response"
	"github.com/samirhembrom/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	target := req.RequestLine.RequestTarget
	headers := headers.NewHeaders()
	if target == "/yourproblem" {
		err := w.WriteStatusLine(response.BadRequest)
		if err != nil {
			return
		}
		err = w.WriteHeaders(headers)
		if err != nil {
			return
		}
		_, err = w.WriteBody([]byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`))
		if err != nil {
			return
		}
		return
	} else if target == "/myproblem" {
		err := w.WriteStatusLine(response.InternalServerError)
		if err != nil {
			return
		}
		err = w.WriteHeaders(headers)
		if err != nil {
			return
		}
		_, err = w.WriteBody([]byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`))
		if err != nil {
			return
		}
		return
	} else if strings.HasPrefix(target, "/httpbin/") {
		endPoint := "https://httpbin.org/" + strings.TrimPrefix(target, "/httpbin/")
		req, err := http.Get(endPoint)
		if err != nil {
			return
		}
		defer req.Body.Close()

		err = w.WriteStatusLine(response.StatusCode(req.StatusCode))
		if err != nil {
			return
		}
		reqContType := req.Header.Get("Content-Type")
		if len(reqContType) == 0 {
			reqContType = "text/plain"
		}
		headers.Set("Transfer-Encoding", "chunked")
		headers.Set("Content-Type", reqContType)
		err = w.WriteHeaders(headers)
		if err != nil {
			return
		}

		buf := make([]byte, 1024, 1024)
		for {
			n, err := req.Body.Read(buf)
			if n > 0 {
				if _, err = w.WriteChunkedBody(buf[:n]); err != nil {
					return
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				return
			}

		}
		_, err = w.WriteChunkedBodyDone()
		if err != nil {
			return
		}
		return

	}
	err := w.WriteStatusLine(response.StatusOK)
	if err != nil {
		return
	}
	err = w.WriteHeaders(headers)
	if err != nil {
		return
	}
	_, err = w.WriteBody([]byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`))
	if err != nil {
		return
	}
}
