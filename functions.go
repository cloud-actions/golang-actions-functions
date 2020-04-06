package main

import (
	"fmt"
	"log"
	"net/http"
)

func AddFunctionHandlers(s *HTTPServer) {
	s.router.HandleFunc("/HttpTrigger", s.httpTrigger())
	s.router.HandleFunc("/HttpTrigger2", s.httpTrigger2())
	s.router.HandleFunc("/TimerTrigger", s.timerTrigger())
}

func (s *HTTPServer) httpTrigger() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// we are modifying this variable, so it needs
		// to be within the function
		entity := "Functions"
		if val := r.FormValue("name"); val != "" {
			entity = val
		}
		fmt.Fprintf(w, "Hello %s!\n", entity)

	}
}

func (s *HTTPServer) httpTrigger2() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// the below is primarily for testing. we have amended functions.json
		// so this should never get hit inside the functions runtime.
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		map1 := map[string]interface{}{}
		err := s.decode(w, r, &map1)
		if err != nil {
			errorText := fmt.Sprintf("Error parsing JSON: %s", err)
			log.Printf("%s", err)
			errorJSON := map[string]interface{}{"error": errorText}
			s.respond(w, r, errorJSON, 500)
			return
		}

		if val, ok := map1["name"]; ok {
			if greeting, ok := val.(string); ok {
				map1["greeting"] = fmt.Sprintf("Hello, %s!", greeting)
			}
		}

		s.respond(w, r, map1, 200)
	}
}

func (s *HTTPServer) timerTrigger() http.HandlerFunc {
	type TimerRequest struct {
		Data struct {
			TimerJSON string `json:"myTimer"`
		} `json:"Data"`
		Metadata struct {
			SysJSON string `json:"sys"`
		} `json:"Metadata"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		t := TimerRequest{}
		err := s.decode(w, r, &t)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		log.Printf("TimerTrigger: %s", t.Data.TimerJSON)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{}")
	}
}
