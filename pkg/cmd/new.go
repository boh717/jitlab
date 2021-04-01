package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func NewTicket() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "new",
		Short: "Pick new issue",
		Long:  `Run this command to pick a new jira issue to work on`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Picking new issue...")
		},
	}

	return newCmd

}
