package question

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/boh717/jitlab/pkg/gitlab"
	"github.com/boh717/jitlab/pkg/jira"
)

type QuestionService interface {
	AskForBoard(boards []jira.Board) (jira.Board, error)
	AskForColumns(columns []jira.Column) ([]string, error)
	AskForRepository(repositories []gitlab.Repository) (gitlab.Repository, error)
	AskForIssue(issues []jira.Issue) (jira.Issue, error)
}

type QuestionServiceImpl struct{}

func (q QuestionServiceImpl) AskForBoard(boards []jira.Board) (jira.Board, error) {

	var chosenBoard jira.Board
	var boardNames []string

	for _, value := range boards {
		if value.Location.DisplayName != "" {
			boardNames = append(boardNames, value.Location.DisplayName)
		}
	}

	question := &survey.Select{
		Message: "Which board do you want to track?",
		Options: boardNames,
	}

	answer := ""

	err := survey.AskOne(question, &answer)
	if err != nil {
		return chosenBoard, err
	}

	getChosenBoard(boards, answer, &chosenBoard)

	return chosenBoard, nil

}

func (q QuestionServiceImpl) AskForColumns(columns []jira.Column) ([]string, error) {

	var columnNames []string

	for _, value := range columns {
		columnNames = append(columnNames, value.Name)
	}

	question := &survey.MultiSelect{
		Message: "Which columns do you want to read from?",
		Options: columnNames,
	}

	answer := []string{}

	err := survey.AskOne(question, &answer)
	if err != nil {
		return nil, err
	}

	return answer, nil

}

func (q QuestionServiceImpl) AskForRepository(repositories []gitlab.Repository) (gitlab.Repository, error) {

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

func (q QuestionServiceImpl) AskForIssue(issues []jira.Issue) (jira.Issue, error) {

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

func getChosenRepo(repos []gitlab.Repository, repoName string, chosenRepo *gitlab.Repository) {

	for _, v := range repos {
		if v.Name == repoName {
			*chosenRepo = v
		}
	}
}

func getChosenBoard(boards []jira.Board, boardName string, chosenBoard *jira.Board) {

	for _, v := range boards {
		if v.Location.DisplayName == boardName {
			*chosenBoard = v
		}
	}
}

func getChosenIssue(repos []jira.Issue, issueName string, chosenIssue *jira.Issue) {

	for _, v := range repos {
		if v.Fields.Summary == issueName {
			*chosenIssue = v
		}
	}
}
