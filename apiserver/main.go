package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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

func main() {

	keepalive := true
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { handlerRoot(w, r, keepalive) })
	http.HandleFunc("/v1/hello", func(w http.ResponseWriter, r *http.Request) { handlerHello(w, r, keepalive) })

	addr := ":3000"

	log.Printf("serving HTTP on TCP %s", addr)

	if err := listenAndServe(addr, nil, keepalive); err != nil {
		log.Fatalf("listenAndServe: %s: %v", addr, err)
	}
}

func listenAndServe(addr string, handler http.Handler, keepalive bool) error {
	server := &http.Server{Addr: addr, Handler: handler}
	server.SetKeepAlivesEnabled(keepalive)
	return server.ListenAndServe()
}

func handlerRoot(w http.ResponseWriter, r *http.Request, keepalive bool) {
	msg := fmt.Sprintf("handlerRoot: url=%s from=%s", r.URL.Path, r.RemoteAddr)
	log.Print(msg)

	if r.URL.Path != "/" {
		io.WriteString(w, header)
		io.WriteString(w, fmt.Sprintf("<h2>path not found!</h2>path not found: [%s]", r.URL.Path))
		io.WriteString(w, footer)
		return
	}

	io.WriteString(w, header)
	io.WriteString(w, "<h2>root handler</h2>nothing to see here")
	io.WriteString(w, footer)
}

func handlerHello(w http.ResponseWriter, r *http.Request, keepalive bool) {
	msg := fmt.Sprintf("handlerHello: url=%s from=%s", r.URL.Path, r.RemoteAddr)
	log.Print(msg)

	io.WriteString(w, header)
	io.WriteString(w, "<h2>hello handler</h2>hello world")
	io.WriteString(w, footer)
}
