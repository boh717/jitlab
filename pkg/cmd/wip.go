package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func WipTicket() *cobra.Command {
	wipCmd := &cobra.Command{
		Use:   "wip",
		Short: "Pick WIP issue",
		Long:  `Run this command to continue to work on a jira ticket (read from your WIP columns). This is useful when you have to work on multiple subtasks (e.g. story) in different folders (i.e. microservices).`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Picking WIP issue...")
		},
	}

	return wipCmd

}
