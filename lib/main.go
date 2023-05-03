package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

type PullRequestOutput struct {
	PullRequests []PullRequest `json:"values"`
	// Size int `json:"size"`
}

type PullRequest struct {
	Links  Links  `json:"links"`
	Title  string `json:"title"`
	Author struct {
		DisplayName string `json:"display_name"`
	} `json:"author"`
}

type GetRepositoriesOutput struct {
	Repositories []Repository `json:"values"`
}

type Repository struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
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
	Workspace              string
	ProjectKeys            []string
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

	if input.ProjectKeys != nil {
		for _, projectKey := range input.ProjectKeys {
			repos, err := GetRepositoriesByProject(projectKey, input.Workspace, input)
			if err != nil {
				return err
			}

			for _, repo := range repos {
				repoPrs, err := GetPullRequestsByRepo(input.Workspace, repo, input)
				if err != nil {
					return err
				}
				prs = append(prs, repoPrs...)
			}
		}
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

func GetPullRequestsByRepo(workspace string, repo Repository, input *NotifySlackInput) ([]PullRequest, error) {
	repos := []PullRequest{}
	pageLength := 50
	curPage := 1

	for {
		client := http.Client{Timeout: 5 * time.Second}

		url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/pullrequests", workspace, repo.Name)
		req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
		if err != nil {
			return nil, fmt.Errorf("Error creating get repository pull requests url: %s", err.Error())
		}

		q := req.URL.Query()
		q.Set("pagelen", strconv.Itoa(pageLength))
		q.Set("page", strconv.Itoa(curPage))
		req.URL.RawQuery = q.Encode()

		fmt.Println(req.URL.String())

		req.SetBasicAuth(input.BitbucketAppPassUser, input.BitbucketAppPassSecret)

		res, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving pull requests by repository: %s", err.Error())
		}

		defer res.Body.Close()

		pullRequestOutput := &PullRequestOutput{}
		err = json.NewDecoder(res.Body).Decode(pullRequestOutput)
		if err != nil {
			return nil, fmt.Errorf("Error decoding response output: %s", err.Error())
		}

		if len(pullRequestOutput.PullRequests) == 0 {
			break
		}
		repos = append(repos, pullRequestOutput.PullRequests...)
		curPage++
	}
	return repos, nil
}

func GetRepositoriesByProject(projectKey string, workspace string, input *NotifySlackInput) ([]Repository, error) {
	repos := []Repository{}
	pageLength := 50
	curPage := 1

	for {
		client := http.Client{Timeout: 5 * time.Second}

		url := fmt.Sprintf(`https://api.bitbucket.org/2.0/repositories/%s`, workspace)
		req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
		if err != nil {
			return nil, fmt.Errorf("Error creating repositories requests url: %s", err.Error())
		}

		q := req.URL.Query()
		q.Add("q", fmt.Sprintf(`project.key="%s"`, projectKey))
		q.Set("pagelen", strconv.Itoa(pageLength))
		q.Set("page", strconv.Itoa(curPage))
		req.URL.RawQuery = q.Encode()

		fmt.Println(req.URL.String())

		req.SetBasicAuth(input.BitbucketAppPassUser, input.BitbucketAppPassSecret)

		res, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving repositories by project key: %s", err.Error())
		}

		defer res.Body.Close()

		getRepositoriesOutput := &GetRepositoriesOutput{}
		err = json.NewDecoder(res.Body).Decode(getRepositoriesOutput)
		if err != nil {
			return nil, fmt.Errorf("Error decoding response output: %s", err.Error())
		}

		if len(getRepositoriesOutput.Repositories) == 0 {
			break
		}
		repos = append(repos, getRepositoriesOutput.Repositories...)
		curPage++
	}
	return repos, nil
}

func postToSlack(prs []PullRequest, input *NotifySlackInput) error {

	slackClient := slack.New(
		input.SlackToken,
		slack.OptionDebug(input.Debug),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)

	msg := "*Open PR Digest*\n"

	for _, pr := range prs {
		msg += fmt.Sprintf("%s - %s - %s\n", pr.Author.DisplayName, pr.Title, pr.Links.Html.Href)
	}
	// TODO add link title instead of full url
	_, _, err := slackClient.PostMessage(input.SlackChannel, slack.MsgOptionText(msg, true))
	if err != nil {
		return err
	}
	return nil
}
