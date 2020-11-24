package main

import (
	"encoding/json"
	"io/ioutil"
)

type Account struct {
	UserId            string `json:"user_id"`
	GithubAccessToken string `json:"github_access_token"`
	Subscribed        bool   `json:"subscribed"`
	LastActiveId      string `json:"last_active_id"`
}

func save(data Account) error {
	bytes, _ := json.Marshal(data)
	return ioutil.WriteFile("data/"+data.UserId+".json", bytes, 0644)
}

func getAllAccounts() []Account {
	response := make([]Account, 0)
	files, _ := ioutil.ReadDir("data")
	for _, f := range files {
		if f.Name() == "README.md" {
			continue
		}
		bytes, err := ioutil.ReadFile("data/" + f.Name())
		if err != nil {
			continue
		}
		data := new(Account)
		_ = json.Unmarshal(bytes, data)
		response = append(response, *data)
	}
	return response
}
