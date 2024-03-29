package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Message represents the JSON request structure expected from the cat_http script
type Message struct {
	Response string `json:"response"`
}

type Request struct {
	Model       string `json:"model"`
	Prompt      string `json:"prompt"`
  Stream      bool `json:"stream"` 
}

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Server listening on port 11434")
	http.ListenAndServe(":11434", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var msg Message
	var request Request

	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

  fmt.Println("Prompt :", request.Prompt)
  fmt.Println("Stream :", request.Stream)
  fmt.Println("Model :", request.Model)

  msg.Response = request.Prompt

	// Create the response message
	response := Message{
		Response: msg.Response, // Echoing back the received content
	}

	// Encode the response to JSON
	respJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	// Set Content-Type and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respJSON)
}

