package cmd

import (
	"fmt"
	"os"

	lib "github.com/chrispruitt/pr-notifier/lib"

	"github.com/spf13/cobra"
)

var (
	notifySlackInput lib.NotifySlackInput
)

func init() {
	NotifyCmd.PersistentFlags().StringVarP(&notifySlackInput.BitbucketAppPassUser, "user", "u", "", "bitbucket app username")
	NotifyCmd.Flags().StringArrayVarP(&notifySlackInput.Usernames, "authors", "a", nil, "bitbucket usernames, UUIDs, or authors of PRs")
	NotifyCmd.Flags().StringVarP(&notifySlackInput.Workspace, "workspace", "w", "", "bitbucket workspace key. Required when --project-key is set.")
	NotifyCmd.Flags().StringArrayVar(&notifySlackInput.ProjectKeys, "project-key", nil, "bitbucket workspace key. Required when --project-key is set.")
	NotifyCmd.Flags().StringVarP(&notifySlackInput.BitbucketAppPassSecret, "password", "p", "", "bitbucket app password")
	NotifyCmd.Flags().StringVarP(&notifySlackInput.SlackChannel, "channel", "c", "", "slack channel")
	NotifyCmd.Flags().StringVarP(&notifySlackInput.SlackToken, "token", "t", "", "slack token")
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
