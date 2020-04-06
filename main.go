package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
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
	// set by --port
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i-1] == "--port" {
			listenAddr = ":" + os.Args[i]
		}
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
	s.router.HandleFunc("/animal", s.httpIndexWithParam("Python"))
	s.router.HandleFunc("/api", s.httpHasAPIVersion(s.httpAPI()))
	s.router.HandleFunc("/greeting", s.httpGreeting())
	s.router.HandleFunc("/echo", s.httpEcho())
	s.router.HandleFunc("/host", s.httpHost())
	s.router.HandleFunc("/template", s.httpTemplate("index.html"))
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

func (s *HTTPServer) httpIndex() http.HandlerFunc {
	animal := "Gopher"
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", animal)
	}
}

func (s *HTTPServer) httpAPI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiVersion := r.FormValue("api-version")
		fmt.Fprintf(w, "api-version: %s\n", apiVersion)
	}
}

func (s *HTTPServer) httpIndexWithParam(animal string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", animal)
	}
}

func (s *HTTPServer) httpGreeting() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
	}
	type response struct {
		Greeting string `json:"greeting,omitempty"`
		Error    string `json:"error,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			s.respond(w, r, &response{Error: "Method should be POST"}, http.StatusMethodNotAllowed)
			return
		}
		req := &request{}
		err := s.decode(w, r, &req)
		if err != nil {
			s.respond(w, r, &response{Error: err.Error()}, 500)
			return
		}
		res := &response{
			Greeting: fmt.Sprintf("Hello %s!", req.Name),
		}
		s.respond(w, r, res, 200)
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

func (s *HTTPServer) httpHost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if hostname, err := os.Hostname(); err == nil {
			fmt.Fprintf(w, "Hostname: %s\n", hostname)
		}
	}
}

func (s *HTTPServer) httpTemplate(files ...string) http.HandlerFunc {
	var (
		init   sync.Once
		tpl    *template.Template
		tplerr error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tpl, tplerr = template.ParseFiles(files...)
		})
		if tplerr != nil {
			http.Error(w, tplerr.Error(), http.StatusInternalServerError)
			return
		}
		// use template
		items := []string{"the", "quick", "brown", "fox", "jumped"}
		data := struct {
			Title  string
			Items  []string
			Result string
		}{
			Title:  "Title",
			Items:  items,
			Result: strings.Join(items, " "),
		}
		if err := tpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
}

func (s *HTTPServer) httpHasAPIVersion(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("api-version") == "" {
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}
}

func (s *HTTPServer) httpLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}
