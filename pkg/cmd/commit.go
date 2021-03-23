package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Commits() *cobra.Command {
	commitCmd := &cobra.Command{
		Use:   "commit",
		Short: "Commit your changes",
		Long:  `Commit your changes with a commit message following a pattern (for example, you may want to include a Jira ticket reference)`,
		Run: func(cmd *cobra.Command, args []string) {
			message, _ := cmd.Flags().GetString("message")
			fmt.Printf("Here's your commit message %s", message)
		},
	}

	requiredFlag := ""
	commitCmd.Flags().StringVarP(&requiredFlag, "message", "m", "", "Your commit message")
	commitCmd.MarkFlagRequired("message")

	return commitCmd

}
