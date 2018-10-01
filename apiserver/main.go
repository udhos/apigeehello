package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	header = `<!DOCTYPE html>
<html>
  <head>
    <title>gowebhello root page</title>
  </head>
  <body>
`

	footer = `</body>
</html>
`
)

type handlerFunc func(w http.ResponseWriter, r *http.Request)

var html bool

func main() {

	keepalive := true

	register("/", func(w http.ResponseWriter, r *http.Request) { handlerRoot(w, r, keepalive) })
	register("/v1/hello", func(w http.ResponseWriter, r *http.Request) { handlerHello(w, r, keepalive) })

	addr := ":3000"

	if os.Getenv("HTML") != "" {
		html = true
	}

	log.Printf("serving HTTP on TCP %s html=%v HTML=[%s]", addr, html, os.Getenv("HTML"))

	if err := listenAndServe(addr, nil, keepalive); err != nil {
		log.Fatalf("listenAndServe: %s: %v", addr, err)
	}
}

func register(path string, handler handlerFunc) {
	log.Printf("registering path: [%s]", path)
	http.HandleFunc(path, handler)
}

func listenAndServe(addr string, handler http.Handler, keepalive bool) error {
	server := &http.Server{Addr: addr, Handler: handler}
	server.SetKeepAlivesEnabled(keepalive)
	return server.ListenAndServe()
}

func sendHeader(w http.ResponseWriter) {
	if html {
		io.WriteString(w, header)
	}
}

func sendFooter(w http.ResponseWriter) {
	if html {
		io.WriteString(w, footer)
	}
}

func sendTag(w http.ResponseWriter, tag, text string) {
	if html {
		io.WriteString(w, "<"+tag+">")
	}
	io.WriteString(w, text)
	if html {
		io.WriteString(w, "</"+tag+">")
	}
}

func handlerRoot(w http.ResponseWriter, r *http.Request, keepalive bool) {
	msg := fmt.Sprintf("handlerRoot: url=%s from=%s", r.URL.Path, r.RemoteAddr)
	log.Print(msg)

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		sendHeader(w)
		sendTag(w, "h2", "path not found!\n")
		io.WriteString(w, fmt.Sprintf("path not found: [%s]\n", r.URL.Path))
		sendFooter(w)
		return
	}

	sendHeader(w)
	sendTag(w, "h2", "root handler\n")
	io.WriteString(w, "nothing to see here\n")
	sendFooter(w)
}

func handlerHello(w http.ResponseWriter, r *http.Request, keepalive bool) {
	msg := fmt.Sprintf("handlerHello: url=%s from=%s", r.URL.Path, r.RemoteAddr)
	log.Print(msg)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	sendHeader(w)
	sendTag(w, "h2", "hello handler\n")
	io.WriteString(w, "hello world\n")
	sendFooter(w)
}
