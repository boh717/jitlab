package question

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/boh717/jitlab/pkg/jira"
)

func AskForBoard(boards []jira.Board) (jira.Board, error) {

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

func AskForColumns(columns []jira.Column) ([]string, error) {

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

func getChosenBoard(boards []jira.Board, boardName string, chosenBoard *jira.Board) {

	for _, v := range boards {
		if v.Location.DisplayName == boardName {
			*chosenBoard = v
		}
	}
}
