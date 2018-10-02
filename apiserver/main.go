package main

import (
	"encoding/json"
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
    <title>apiserver root page</title>
  </head>
  <body>
`

	footer = `</body>
</html>
`
)

type responseHello struct {
	Message string `json:"message"`
	Age     int    `json:"age"`
}

type responseError struct {
	Message string `json:"message"`
}

type handlerFunc func(w http.ResponseWriter, r *http.Request)

var html bool

func main() {

	keepalive := true

	register("/", func(w http.ResponseWriter, r *http.Request) { handlerRoot(w, r, keepalive, "/") })
	register("/v1/hello", func(w http.ResponseWriter, r *http.Request) { handlerHello(w, r, keepalive, "/v1/hello") })
	register("/v1/hello/", func(w http.ResponseWriter, r *http.Request) { handlerHello(w, r, keepalive, "/v1/hello/") })

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

func sendNotFound(label string, w http.ResponseWriter, r *http.Request, useJson bool) {
	msg := fmt.Sprintf("%s: url=%s from=%s json=%v - PATH NOT FOUND", label, r.URL.Path, r.RemoteAddr, useJson)
	log.Print(msg)

	notFound := fmt.Sprintf("path not found: [%s]", r.URL.Path)

	w.WriteHeader(http.StatusNotFound)

	if useJson {
		resp := responseError{notFound}
		b, errMarshal := json.Marshal(resp)
		if errMarshal != nil {
			log.Printf("json marshal: %v", errMarshal)
			return
		}
		w.Write(b)
		io.WriteString(w, "\n")
		return
	}

	sendHeader(w)
	sendTag(w, "h2", "path not found\n")
	io.WriteString(w, notFound+"\n")
	sendFooter(w)
}

func acceptJson(r *http.Request) bool {
	var found bool

	for k, v := range r.Header {
		if k == "Accept" {
			for _, vv := range v {
				if vv == "application/json" {
					found = true
				}
			}
		}
	}

	return found
}

func handlerRoot(w http.ResponseWriter, r *http.Request, keepalive bool, path string) {

	useJson := acceptJson(r)

	if r.URL.Path != path {
		sendNotFound("handlerRoot", w, r, useJson)
		return
	}

	msg := fmt.Sprintf("handlerRoot: url=%s from=%s json=%v", r.URL.Path, r.RemoteAddr, useJson)
	log.Print(msg)

	nothing := fmt.Sprintf("nothing to see here: [%s]", r.URL.Path)

	if useJson {
		resp := responseError{nothing}
		b, errMarshal := json.Marshal(resp)
		if errMarshal != nil {
			log.Printf("json marshal: %v", errMarshal)
			return
		}
		w.Write(b)
		io.WriteString(w, "\n")
		return
	}

	sendHeader(w)
	sendTag(w, "h2", "root handler\n")
	io.WriteString(w, nothing+"\n")
	sendFooter(w)
}

func handlerHello(w http.ResponseWriter, r *http.Request, keepalive bool, path string) {

	useJson := acceptJson(r)

	if r.URL.Path != path {
		sendNotFound("handlerHello", w, r, useJson)
		return
	}

	msg := fmt.Sprintf("handlerHello: url=%s from=%s json=%v", r.URL.Path, r.RemoteAddr, useJson)
	log.Print(msg)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	hello := "hello world"

	if useJson {
		resp := responseHello{Message: hello, Age: 17}
		b, errMarshal := json.Marshal(resp)
		if errMarshal != nil {
			log.Printf("json marshal: %v", errMarshal)
			return
		}
		w.Write(b)
		io.WriteString(w, "\n")
		return
	}

	sendHeader(w)
	sendTag(w, "h2", "hello handler\n")
	io.WriteString(w, hello+"\n")
	sendFooter(w)
}
