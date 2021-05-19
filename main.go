package main

import (
	"encoding/json"
	"fmt"
	"github.com/antage/eventsource"
	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"strconv"
	"time"
)

func addUserHandler(writer http.ResponseWriter, request *http.Request) {
	userName := request.FormValue("name")
	sendMessage("", fmt.Sprintf("%s joined", userName))
}

type Message struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

var msgCh chan Message

func postMessageHandler(writer http.ResponseWriter, request *http.Request) {
	msg := request.FormValue("msg")
	name := request.FormValue("name")
	sendMessage(name, msg)
	log.Println("postMessageHandler", msg, name)
}

func sendMessage(name, msg string) {
	// broadcast message
	msgCh <- Message{name, msg}
}

func processMsgCh(es eventsource.EventSource) {
	for msg := range msgCh {
		data, _ := json.Marshal(msg)
		es.SendEventMessage(string(data), "", strconv.Itoa(time.Now().Nanosecond()))

	}
}

func main() {
	msgCh = make(chan Message)
	es := eventsource.New(nil, nil)
	defer es.Close()

	go processMsgCh(es)
	
	mux := pat.New()
	mux.Post("/messages", postMessageHandler)
	mux.Handle("/stream", es)
	mux.Post("/users", addUserHandler)
	mux.Delete("/users", leftUserHandler)

	neg := negroni.Classic()
	neg.UseHandler(mux)

	http.ListenAndServe(":3000", neg)
}

func leftUserHandler(writer http.ResponseWriter, request *http.Request) {
	userName := request.FormValue("name")
	sendMessage("", fmt.Sprintf("%s left", userName))
}
