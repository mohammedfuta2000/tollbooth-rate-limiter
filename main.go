package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/didip/tollbooth"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func endpointHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("testing endpoint")
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	message := Message{
		Status: "Successful",
		Body:   "Hi, you have reached the API. How may I help you?",
	}
	err := json.NewEncoder(writer).Encode(&message)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	message := Message{
		Status: "Request failed",
		Body:   "the API is at capacity",
	}
	jsonMessage, _ := json.Marshal(message)
	tbLimiter := tollbooth.NewLimiter(1, nil)
	tbLimiter.SetMessageContentType("application/json")
	tbLimiter.SetMessage(string(jsonMessage))

	http.Handle("/ping", tollbooth.LimitFuncHandler(tbLimiter, endpointHandler))
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Println(err)
	}
}
