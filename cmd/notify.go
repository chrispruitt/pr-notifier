package cmd

import (
	lib "github.com/chrispruitt/pr-notifier/lib"

	"github.com/spf13/cobra"
)

var (
	notifySlackInput lib.NotifySlackInput
)

func init() {
	NotifyCmd.PersistentFlags().StringVarP(&notifySlackInput.BitbucketAppPassUser, "user", "u", "", "bitbucket app username")
	NotifyCmd.Flags().StringArrayVarP(&notifySlackInput.Usernames, "authors", "a", nil, "bitbucket usernames, UUIDs, or authors of PRs")
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
		err := lib.NotifySlack(&notifySlackInput)
		if err != nil {
			panic(err)
		}
	},
}
