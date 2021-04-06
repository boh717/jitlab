package question

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/boh717/jitlab/pkg/gitlab"
)

func AskForRepository(repositories []gitlab.Repository) (gitlab.Repository, error) {

	var chosenRepo gitlab.Repository
	var repoNames []string

	for _, value := range repositories {
		repoNames = append(repoNames, value.Name)
	}

	question := &survey.Select{
		Message: "Which repository are you looking for?",
		Options: repoNames,
	}

	answer := ""

	err := survey.AskOne(question, &answer)
	if err != nil {
		return chosenRepo, err
	}

	getChosenRepo(repositories, answer, &chosenRepo)

	return chosenRepo, nil

}

func getChosenRepo(repos []gitlab.Repository, repoName string, chosenRepo *gitlab.Repository) {

	for _, v := range repos {
		if v.Name == repoName {
			*chosenRepo = v
		}
	}
}
