package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/boh717/jitlab/pkg/gitlab"
	"github.com/spf13/cobra"
)

func MergeRequest() *cobra.Command {
	mrCmd := &cobra.Command{
		Use:   "mr",
		Short: "Create a new merge request",
		Long:  `Run this command to create a new merge request using target branch and (limited) options of your choice`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Creating new merge request...")
			targetBranch, _ := cmd.Flags().GetString("target-branch")
			removeSourceBranch, _ := cmd.Flags().GetBool("remove-source-branch")
			squash, _ := cmd.Flags().GetBool("squash")

			branch, err := gitService.GetCurrentBranch()
			if err != nil {
				log.Fatalln("Error getting branch", err)
			}

			pushErr := gitService.Push(branch)
			if pushErr != nil {
				log.Fatalln("Error pushing branch", pushErr)
			}

			var currentRepository gitlab.Repository

			file, err := ioutil.ReadFile(".repo")
			if err != nil {
				log.Fatalln("Error reading repository file \".repo\"", err)
			}

			json.Unmarshal(file, &currentRepository)
			projectId := fmt.Sprintf("%d", currentRepository.ID)

			title, err := gitService.CreateTitleFromBranch(branch)
			if err != nil {
				log.Fatalln("Error creating title from branch", err)
			}

			resp, err := gitlabService.CreateMergeRequest(projectId, branch, targetBranch, title, removeSourceBranch, squash)
			if err != nil {
				log.Fatalln("Error creating merge request", err)
			}
			log.Printf("Merge request created: %s", resp.Url)

		},
	}

	var targetBranch string
	var removeSourceBranch bool
	var squash bool

	mrCmd.Flags().StringVar(&targetBranch, "target-branch", "master", "Target branch for merge request")
	mrCmd.Flags().BoolVar(&removeSourceBranch, "remove-source-branch", true, "Remove source branch when merging")
	mrCmd.Flags().BoolVar(&squash, "squash", true, "Squash commits when merging")

	return mrCmd

}
