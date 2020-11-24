package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	SlackBotToken      = flag.String("slack_bot_token", "", "Bot token for slack application")
	GithubClientId     = flag.String("github_client_id", "", "client id of github application")
	GithubClientSecret = flag.String("github_client_secret", "", "client secret for github application")
	Address            = flag.String("address", "", "public address of the application")
)

func main() {
	flag.Parse()
	router := mux.NewRouter().StrictSlash(true)
	StartJobToSendNotification()
	RegisterMessageHandler(router)
	err := http.ListenAndServe(":8888", router)
	log.Print(err)
	log.Print("server is closed")
}
