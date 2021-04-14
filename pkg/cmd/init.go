package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/boh717/jitlab/pkg/question"
	"github.com/spf13/cobra"
)

func InitRepo() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Configure your repository",
		Long:  `Run this command in every git repo you want to use Jitlab`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Init repo...")

			currentPath, err := os.Getwd()
			if err != nil {
				log.Fatalln(err)
			}
			currentDir := path.Base(currentPath)

			repositories, err := gitlabClient.SearchProject(currentDir)
			if err != nil {
				log.Fatalln(err)
			}

			if len(repositories) == 0 {
				log.Fatalf("Your search \"%s\" didn't match any project", currentDir)
			}

			var chosenRepo = repositories[0]
			if len(repositories) > 1 {
				chosenRepo, err = question.AskForRepository(repositories)
				if err != nil {
					log.Fatalln(err)
				}
			}

			file, _ := json.MarshalIndent(chosenRepo, "", " ")

			if err := ioutil.WriteFile(".repo", file, 0644); err != nil {
				log.Fatalln(err)
			}

		},
	}

	return initCmd
}
