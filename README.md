# NotifyMe
Have you ever face the problem where you miss important github 
notification because of all the noise between subscribed events. 
Like someone's mention, review request etc. I faced this challenge, 
and this app is to address this issue. It integrated with the slack 
and sends all the important notifications.

---
# How to run

First build the application via
`go build .`

Then you can run the app by `./notifyme --slack_bot_token=slacktoken --github_client_id=clientid --github_client_secret=secret --address=applicationAddress`

---
# Scope Of Improvement
There are many things that needs to be improved. 
1. Message format.
2. Integration with other tools having important notifications like discuss.

---