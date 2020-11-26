package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const slackNotifPath = "https://slack.com/api/chat.postMessage?token=%s&channel=%s&text=%s&blocks=%s&pretty=1"

type GithubResponse struct {
	Id         string           `json:"id"`
	Reason     string           `json:"reason"`
	Subject    GithubSubject    `json:"subject"`
	UpdatedAt  string           `json:"updated_at"`
	Repository GithubRepository `json:"repository"`
}

type GithubSubject struct {
	Title string `json:"title"`
	Url   string `json:"url"`
	Type  string `json:"type"`
}

type GithubRepository struct {
	Name    string `json:"name"`
	HtmlUrl string `json:"html_url"`
}

func StartJobToSendNotification() {
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			data := getAllAccounts()
			if len(data) == 0 {
				log.Print("no users to send notification to. Hence skipping")
				continue
			}
			for _, d := range data {
				if !d.Subscribed {
					continue
				}

				lastUpdatedAt, err := checkAndSendNotificationIfRequired(d)
				if err != nil {
					log.Print(err)
					continue
				}

				d.LastActiveId = lastUpdatedAt
				log.Print(save(d))
			}
		}
	}()
}

//TODO make github call paginated
func checkAndSendNotificationIfRequired(data Account) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/notifications?all=true", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+data.GithubAccessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	do, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	response := readResponse(do)
	if len(response) == 0 {
		return "", nil
	}

	if len(data.LastActiveId) == 0 {
		return response[0].UpdatedAt, nil
	}
	notifToSend := make([]*GithubResponse, 0)
	for _, r := range response {
		if r.UpdatedAt > data.LastActiveId && r.Reason != "subscribed" {
			notifToSend = append(notifToSend, r)
		}
	}
	sendSlackNotificationForGithub(data, notifToSend)
	return response[0].UpdatedAt, nil
}

func sendSlackNotificationForGithub(d Account, notifFor []*GithubResponse) {
	if len(notifFor) == 0 {
		return
	}

	for _, n := range notifFor {
		slackPath := fmt.Sprintf(slackNotifPath, *SlackBotToken, d.UserId,
			url.PathEscape("You have new notification on github."),
			url.PathEscape(fmt.Sprintf(`[
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "This is related to repository %s because of %s over %s"
			}
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "Reason: *%s*"
			},
			"accessory": {
				"type": "button",
				"text": {
					"type": "plain_text",
					"text": "Take me there",
					"emoji": true
				},
				"url": "%s",
			}
		},
		{
			"type": "divider"
		}
	]`, n.Repository.Name, n.Reason, n.Subject.Type, n.Reason, "https://github.com/"+strings.Trim(n.Subject.Url, "https://api.github.com/repos/"))))
		req, err := http.NewRequest("POST", slackPath, nil)
		if err != nil {
			log.Print(err)
			return
		}
		req.Header.Set("accept", "application/json")
		_, err = http.DefaultClient.Do(req)
		log.Print(err)
	}
}

func readResponse(r *http.Response) []*GithubResponse {
	defer func() { _ = r.Body.Close() }()

	bytes, _ := ioutil.ReadAll(r.Body)
	response := make([]*GithubResponse, 0)
	_ = json.Unmarshal(bytes, &response)
	return response
}
