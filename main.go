package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%s\n", err)
	}
}

func run() error {
	listenAddr := ""
	if val := os.Getenv("LISTEN_ADDR"); val != "" {
		listenAddr = val
	}
	if val := os.Getenv("LISTEN_PORT"); val != "" {
		listenAddr = ":" + val
	}
	if val := os.Getenv("FUNCTIONS_HTTPWORKER_PORT"); val != "" {
		listenAddr = ":" + val
	}
	srv := NewHTTPServer(listenAddr)
	srv.serverName = "hello-gopher"
	AddFunctionHandlers(srv)
	return srv.Start()
}

type HTTPServer struct {
	serverName string
	addr       string
	router     *http.ServeMux
}

func NewHTTPServer(addr string) *HTTPServer {
	s := &HTTPServer{
		serverName: "default",
		addr:       ":80",
		router:     http.DefaultServeMux,
	}
	if addr != "" {
		s.addr = addr
	}
	s.Routes()
	return s
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *HTTPServer) Start() error {
	fmt.Printf("Server \"%s\" listening on %s\n", s.serverName, s.addr)
	return http.ListenAndServe(s.addr, s.httpLog(s.router))
}

func (s *HTTPServer) Routes() {
	s.router.HandleFunc("/", s.httpEcho())
}

func (s *HTTPServer) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (s *HTTPServer) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		// TODO: handle error better
		log.Printf("%s\n", err)
	}
}

func (s *HTTPServer) httpEcho() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Printf("Error: %s\n", b)
			http.Error(w, err.Error(), 500)
			return
		}
		log.Printf("%s\n", b)
		fmt.Fprintf(w, "%s", b)
	}
}

func (s *HTTPServer) httpIndexWithParam(animal string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", animal)
	}
}

func (s *HTTPServer) httpLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}
