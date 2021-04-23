package cmd

import (
	"github.com/spf13/cobra"
)

func Commits() *cobra.Command {
	commitCmd := &cobra.Command{
		Use:   "commit",
		Short: "Commit your changes",
		Long:  `Commit your changes with a commit message following a pattern (for example, you may want to include a Jira ticket reference)`,
		Run: func(cmd *cobra.Command, args []string) {
			commitMessage, _ := cmd.Flags().GetString("message")
			gitService.Commit(commitMessage)
		},
	}

	requiredFlag := ""
	commitCmd.Flags().StringVarP(&requiredFlag, "message", "m", "", "Your commit message")
	commitCmd.MarkFlagRequired("message")

	return commitCmd

}
