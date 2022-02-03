package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type GuessTheAgeRequest struct {
	Name string `json:"name"`
}

type GuessTheAgeResponse struct {
	Age int `json:"age"`
}

type AgifyResponse struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Count int    `json:"count"`
}

func handleGuessTheAge(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePost(w, r)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var body GuessTheAgeRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		if (err.Error() == "EOF") {
			http.Error(w, "Body Required", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pathTemplate := "https://api.agify.io/?name=%s"
	agifyResponse, err := http.Get(fmt.Sprintf(pathTemplate, body.Name))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	agifyResponseData, err := ioutil.ReadAll(agifyResponse.Body)

	var agifyResponseParsed AgifyResponse
	err = json.Unmarshal(agifyResponseData, &agifyResponseParsed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(GuessTheAgeResponse{Age: agifyResponseParsed.Age})
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	_, _ = w.Write(jsonResp)
	return
}
