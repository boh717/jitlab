package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewTicket() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "new",
		Short: "Pick new issue",
		Long:  `Run this command to pick a new jira issue to work on`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Picking new issue...")
			assignedToMe, _ := cmd.Flags().GetBool("me")

			flowType := viper.GetString("board.type")
			projectKey := viper.GetString("board.location.projectkey")
			columns := viper.GetStringSlice("columns")

			issues, err := jiraService.GetIssues(flowType, projectKey, columns, assignedToMe)
			if err != nil {
				log.Fatalln(err)
			}

			chosenIssue, err := questionService.AskForIssue(issues)
			if err != nil {
				log.Fatalln(err)
			}

			cmdErr := gitService.CreateBranch(chosenIssue)
			if cmdErr != nil {
				log.Fatalln(err)
			}
		},
	}

	var currentUserFlag bool

	newCmd.Flags().BoolVar(&currentUserFlag, "me", false, "Only issues assigned to me")

	return newCmd

}
