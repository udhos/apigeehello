package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
)

const (
	version = "0.0"

	headerTitleBefore = `<!DOCTYPE html>
<html>
  <head>
    <title>`

	headerTitleAfter = `</title>
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

type responseEcho struct {
	Message string `json:"message"`
}

type handlerFunc func(w http.ResponseWriter, r *http.Request)

var (
	html       bool
	errorRate  int
	errorCount int
)

func main() {

	keepalive := true

	log.Printf("apiserver %s runtime %s GOMAXPROCS=%d", version, runtime.Version(), runtime.GOMAXPROCS(0))

	register("/", func(w http.ResponseWriter, r *http.Request) { handlerRoot(w, r, keepalive, "/") })
	register("/v1/hello", func(w http.ResponseWriter, r *http.Request) { handlerHello(w, r, keepalive, "/v1/hello") })
	register("/v1/hello/", func(w http.ResponseWriter, r *http.Request) { handlerHello(w, r, keepalive, "/v1/hello/") })
	register("/v1/echo", func(w http.ResponseWriter, r *http.Request) { handlerEcho(w, r, keepalive, "/v1/echo") })
	register("/v1/echo/", func(w http.ResponseWriter, r *http.Request) { handlerEcho(w, r, keepalive, "/v1/echo/") })

	addr := os.Getenv("LISTEN")

	if addr == "" {
		addr = ":3000"
	}

	if os.Getenv("HTML") != "" {
		html = true
	}

	errorRateStr := os.Getenv("ERROR")
	var errRate error
	errorRate, errRate = strconv.Atoi(errorRateStr)
	if errRate != nil {
		log.Printf("bad error rate: ERROR=[%s]: %v", errorRateStr, errRate)
	}

	log.Printf("errorRate=%d ERROR=[%s]", errorRate, errorRateStr)

	log.Printf("serving HTTP on TCP %s LISTEN=[%s] html=%v HTML=[%s]", addr, os.Getenv("LISTEN"), html, os.Getenv("HTML"))

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

func sendHeader(w http.ResponseWriter, title string) {
	if html {
		io.WriteString(w, headerTitleBefore)
	}
	io.WriteString(w, title)
	if html {
		io.WriteString(w, headerTitleAfter)
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
	msg := fmt.Sprintf("%s: method=%s url=%s from=%s json=%v - PATH NOT FOUND", label, r.Method, r.URL.Path, r.RemoteAddr, useJson)
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

	sendHeader(w, label+" - not found\n")
	sendTag(w, "h2", "path not found\n")
	io.WriteString(w, notFound+"\n")
	sendFooter(w)
}

func acceptJson(r *http.Request) bool {
	var found bool

	for k, v := range r.Header {
		if k == "Accept" {
			for _, vv := range v {
				log.Printf("Accept: %s", vv)
				if vv == "application/json" {
					found = true
				}
			}
		}
	}

	return found
}

func forceError(label string, w http.ResponseWriter, r *http.Request) bool {
	if errorRate < 1 {
		return false
	}

	errorCount = (errorCount + 1) % errorRate

	log.Printf("%s: forceError method=%s url=%s from=%s - rate=%d count=%d", label, r.Method, r.URL.Path, r.RemoteAddr, errorRate, errorCount)

	if errorCount != 0 {
		return false
	}

	log.Printf("%s: forceError method=%s url=%s from=%s - forcing error", label, r.Method, r.URL.Path, r.RemoteAddr)
	http.Error(w, "Internal server error", 500)
	return true
}

func handlerRoot(w http.ResponseWriter, r *http.Request, keepalive bool, path string) {

	useJson := acceptJson(r)

	if r.URL.Path != path {
		sendNotFound("handlerRoot", w, r, useJson)
		return
	}

	msg := fmt.Sprintf("handlerRoot: method=%s url=%s from=%s json=%v", r.Method, r.URL.Path, r.RemoteAddr, useJson)
	log.Print(msg)

	if forceError("handlerRoot", w, r) {
		return
	}

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

	sendHeader(w, "api root\n")
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

	msg := fmt.Sprintf("handlerHello: method=%s url=%s from=%s json=%v", r.Method, r.URL.Path, r.RemoteAddr, useJson)
	log.Print(msg)

	if forceError("handlerHello", w, r) {
		return
	}

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

	sendHeader(w, "api hello\n")
	sendTag(w, "h2", "hello handler\n")
	io.WriteString(w, hello+"\n")
	sendFooter(w)
}

func handlerEcho(w http.ResponseWriter, r *http.Request, keepalive bool, path string) {

	useJson := acceptJson(r)

	if r.URL.Path != path {
		sendNotFound("handlerEcho", w, r, useJson)
		return
	}

	if r.Method != http.MethodPost {
		log.Printf("handlerEcho: method=%s url=%s from=%s json=%v - method not supported", r.Method, r.URL.Path, r.RemoteAddr, useJson)
		w.Header().Set("Allow", "POST") // required by 405 error
		http.Error(w, r.Method+" method not supported (only POST is supported)", 405)
		return
	}

	msg := fmt.Sprintf("handlerEcho: method=%s url=%s from=%s json=%v", r.Method, r.URL.Path, r.RemoteAddr, useJson)
	log.Print(msg)

	if forceError("handlerEcho", w, r) {
		return
	}

	body, errBody := ioutil.ReadAll(r.Body)
	if errBody != nil {
		log.Printf("handlerEcho: method=%s url=%s from=%s json=%v - body error: %v", r.Method, r.URL.Path, r.RemoteAddr, useJson, errBody)
		http.Error(w, "Internal server error", 500)
		return
	}

	echo := string(body)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if useJson {
		resp := responseEcho{Message: echo}
		b, errMarshal := json.Marshal(resp)
		if errMarshal != nil {
			log.Printf("json marshal: %v", errMarshal)
			return
		}
		w.Write(b)
		io.WriteString(w, "\n")
		return
	}

	sendHeader(w, "api echo\n")
	sendTag(w, "h2", "echo handler\n")
	io.WriteString(w, echo+"\n")
	sendFooter(w)
}
