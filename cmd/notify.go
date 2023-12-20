package cmd

import (
	"fmt"
	"os"
	"strings"

	lib "github.com/chrispruitt/pr-notifier/lib"

	"github.com/spf13/cobra"
)

var (
	notifySlackInput lib.NotifySlackInput
)

func init() {

	userEnv := os.Getenv("BB_USERNAME")
	appPasswordEnv := os.Getenv("BB_APP_PASSWORD")
	workspaceEnv := os.Getenv("WORKSPACE")
	slackChannelEnv := os.Getenv("SLACK_CHANNEL")
	slackTokenEnv := os.Getenv("SLACK_TOKEN")
	authorsEnv := strings.Split(os.Getenv("AUTHORS"), ",")
	projectKeysEnv := strings.Split(os.Getenv("PROJECT_KEYS"), ",")

	NotifyCmd.PersistentFlags().StringVarP(&notifySlackInput.BitbucketAppPassUser, "user", "u", userEnv, "bitbucket app username")
	NotifyCmd.Flags().StringArrayVarP(&notifySlackInput.Usernames, "authors", "a", authorsEnv, "bitbucket usernames, UUIDs, or authors of PRs")
	NotifyCmd.Flags().StringVarP(&notifySlackInput.Workspace, "workspace", "w", workspaceEnv, "bitbucket workspace key. Required when --project-key is set.")
	NotifyCmd.Flags().StringArrayVar(&notifySlackInput.ProjectKeys, "project-key", projectKeysEnv, "bitbucket workspace key. Required when --project-key is set.")
	NotifyCmd.Flags().StringVarP(&notifySlackInput.BitbucketAppPassSecret, "password", "p", appPasswordEnv, "bitbucket app password")
	NotifyCmd.Flags().StringVarP(&notifySlackInput.SlackChannel, "channel", "c", slackChannelEnv, "slack channel")
	NotifyCmd.Flags().StringVarP(&notifySlackInput.SlackToken, "token", "t", slackTokenEnv, "slack token")
	NotifyCmd.Flags().BoolVar(&notifySlackInput.Debug, "debug", false, "enable verbose logging")
}

// Run Command ./pentaho-cli run
var NotifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "post open PRs to slack",
	Run: func(cmd *cobra.Command, args []string) {

		if notifySlackInput.ProjectKeys != nil || len(notifySlackInput.ProjectKeys) > 0 {
			if notifySlackInput.Workspace == "" {
				fmt.Println("Missing --workspace. Required when --project-key is set.")
				os.Exit(1)
			}
		}

		err := lib.NotifySlack(&notifySlackInput)
		if err != nil {
			panic(err)
		}
	},
}
