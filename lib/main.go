package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/slack-go/slack"
)

type PullRequestOutput struct {
	PullRequests []PullRequest `json:"values"`
	// Size int `json:"size"`
}

type PullRequest struct {
	Links Links  `json:"links"`
	Title string `json:"title"`
}

type Links struct {
	Html Html `json:"html"`
}

type Html struct {
	Href string `json:"href"`
}

type NotifySlackInput struct {
	BitbucketAppPassUser   string
	BitbucketAppPassSecret string
	Usernames              []string
	SlackChannel           string
	SlackToken             string
	Debug                  bool
}

// HelloWorld says hello
func NotifySlack(input *NotifySlackInput) error {
	prs := []PullRequest{}

	for _, user := range input.Usernames {
		userPRs, err := GetPullRequestsByUser(user, input)
		if err != nil {
			return err
		}
		prs = append(prs, userPRs...)
	}

	if len(prs) > 0 {
		err := postToSlack(prs, input)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetPullRequestsByUser(user string, input *NotifySlackInput) ([]PullRequest, error) {
	client := http.Client{Timeout: 5 * time.Second}

	url := fmt.Sprintf("https://api.bitbucket.org/2.0/pullrequests/%s?pagelen=50", user)
	fmt.Println(url)
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("Error creating requests: %s", err.Error())
	}

	req.SetBasicAuth(input.BitbucketAppPassUser, input.BitbucketAppPassSecret)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving pull requests: %s", err.Error())
	}

	defer res.Body.Close()

	pullRequestOutput := &PullRequestOutput{}
	err = json.NewDecoder(res.Body).Decode(pullRequestOutput)
	if err != nil {
		return nil, fmt.Errorf("Error decoding response output: %s", err.Error())
	}

	return pullRequestOutput.PullRequests, nil
}

func postToSlack(prs []PullRequest, input *NotifySlackInput) error {

	slackClient := slack.New(
		input.SlackToken,
		slack.OptionDebug(input.Debug),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)

	msg := "*Open PR Digest*\n"

	for _, pr := range prs {
		msg += fmt.Sprintf("PR - %s - %s\n", pr.Title, pr.Links.Html.Href)
	}
	_, _, err := slackClient.PostMessage(input.SlackChannel, slack.MsgOptionText(msg, true))
	if err != nil {
		return err
	}
	return nil
}
