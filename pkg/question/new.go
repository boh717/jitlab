package question

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/boh717/jitlab/pkg/jira"
)

func AskForIssue(issues []jira.Issue) (jira.Issue, error) {

	var chosenIssue jira.Issue
	var issueNames []string

	for _, value := range issues {
		issueNames = append(issueNames, value.Fields.Summary)
	}

	question := &survey.Select{
		Message: "Which issue do you want to work on?",
		Options: issueNames,
	}

	answer := ""

	err := survey.AskOne(question, &answer)
	if err != nil {
		return chosenIssue, err
	}

	getChosenIssue(issues, answer, &chosenIssue)

	return chosenIssue, nil

}

func getChosenIssue(repos []jira.Issue, issueName string, chosenIssue *jira.Issue) {

	for _, v := range repos {
		if v.Fields.Summary == issueName {
			*chosenIssue = v
		}
	}
}
