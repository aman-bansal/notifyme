package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const message = "{\n\t\"blocks\": [\n\t\t{\n\t\t\t\"type\": \"section\",\n\t\t\t\"text\": {\n\t\t\t\t\"type\": \"mrkdwn\",\n\t\t\t\t\"text\": \"Loging to Github\"\n\t\t\t},\n\t\t\t\"accessory\": {\n\t\t\t\t\"type\": \"button\",\n\t\t\t\t\"text\": {\n\t\t\t\t\t\"type\": \"plain_text\",\n\t\t\t\t\t\"text\": \"Log In\",\n\t\t\t\t\t\"emoji\": true\n\t\t\t\t},\n\t\t\t\t\"url\": \"%s\"\n\t\t\t}\n\t\t}\n\t]\n}"
const githubTokenPath = "https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s&state=%s"
const githubAuthorizePath = "https://github.com/login/oauth/authorize?client_id=%s" +
	"&redirect_uri=%s&allow_signup=false&scope=notifications&state=%s"

func RegisterMessageHandler(router *mux.Router) {
	router.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		log.Print(r.Form)
		url := r.Form.Get("response_url")
		encodeUser := base64.StdEncoding.EncodeToString([]byte(r.Form.Get("user_id")))
		body := fmt.Sprintf(message, fmt.Sprintf(githubAuthorizePath, *GithubClientId, "http://"+*Address+"/oauth/redirect", encodeUser))

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(body)))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		req.Header.Set("accept", "application/json")
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods("POST")

	router.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if len(code) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		stateVar := r.URL.Query().Get("state")
		reqURL := fmt.Sprintf(githubTokenPath, *GithubClientId, *GithubClientSecret, code, stateVar)
		req, err := http.NewRequest(http.MethodPost, reqURL, nil)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		req.Header.Set("accept", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer func() { _ = res.Body.Close() }()
		bts, err := ioutil.ReadAll(res.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		token := new(Token)
		_ = json.Unmarshal(bts, token)
		userId, _ := base64.StdEncoding.DecodeString(stateVar)
		log.Print(save(Account{
			UserId:            string(userId),
			GithubAccessToken: token.AccessToken,
			Subscribed:        true,
			LastActiveId:      "",
		}))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}).Methods("GET")
}

type Token struct {
	AccessToken string `json:"access_token"`
}
