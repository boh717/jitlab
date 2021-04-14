package jira

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var client = http.DefaultClient

type JiraService interface {
	GetBoards() ([]Board, error)
	GetBoardColumns(board Board) ([]Column, error)
	GetIssues(flowType string, projectKey string, columns []string, currentUser bool) ([]Issue, error)
}

type JiraServiceImpl struct {
	BaseURL  string
	Token    string
	Username string
}

type Board struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Location struct {
		ProjectID      int    `json:"projectId"`
		DisplayName    string `json:"displayName"`
		ProjectName    string `json:"projectName"`
		ProjectKey     string `json:"projectKey"`
		ProjectTypeKey string `json:"projectTypeKey"`
		Name           string `json:"name"`
	} `json:"location"`
}

type boardBase struct {
	MaxResults int     `json:"maxResults"`
	StartAt    int     `json:"startAt"`
	Total      int     `json:"total"`
	IsLast     bool    `json:"isLast"`
	Values     []Board `json:"values"`
}

type Column struct {
	Name string `json:"name"`
}

type boardConfig struct {
	ColumnConfig struct {
		Columns []Column `json:"columns"`
	} `json:"columnConfig"`
}

type Issue struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Fields struct {
		Summary string `json:"summary"`
	} `json:"fields"`
}

type issueBase struct {
	Total  int     `json:"total"`
	Issues []Issue `json:"issues"`
}

func (j JiraServiceImpl) GetBoards() ([]Board, error) {
	uri := "/rest/agile/1.0/board"

	req, _ := http.NewRequest(http.MethodGet, j.BaseURL+uri, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", j.Token))

	board := new(boardBase)
	err := doRequest(*req, board)
	if err != nil {
		return nil, err
	}

	return board.Values, nil
}

func (j JiraServiceImpl) GetBoardColumns(board Board) ([]Column, error) {
	uri := fmt.Sprintf("/rest/agile/1.0/board/%d/configuration", board.ID)

	req, _ := http.NewRequest(http.MethodGet, j.BaseURL+uri, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", j.Token))

	boardConfig := new(boardConfig)
	err := doRequest(*req, boardConfig)
	if err != nil {
		return nil, err
	}

	return boardConfig.ColumnConfig.Columns, nil
}

func (j JiraServiceImpl) GetIssues(flowType string, projectKey string, columns []string, currentUser bool) ([]Issue, error) {
	searchString := buildSearchString(flowType, projectKey, columns, currentUser)
	uri := fmt.Sprintf("/rest/api/3/search?jql=%s&fields=summary,assignee&maxResults=50", url.QueryEscape(searchString))

	req, _ := http.NewRequest(http.MethodGet, j.BaseURL+uri, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", j.Token))

	issue := new(issueBase)
	err := doRequest(*req, issue)
	if err != nil {
		return nil, err
	}

	return issue.Issues, nil

}

func buildSearchString(flowType string, projectKey string, columns []string, currentUser bool) string {
	var searchString strings.Builder

	searchString.WriteString(fmt.Sprintf("project = %s AND ", projectKey))
	quotedColumns := "\"" + strings.Join(columns, "\", \"") + "\""
	searchString.WriteString(fmt.Sprintf("status in (%s) AND ", quotedColumns))
	searchString.WriteString("type != Epic ")

	if currentUser {
		searchString.WriteString("AND assignee=currentUser() ")
	}

	if flowType == "kanban" || flowType == "simple" {
		searchString.WriteString("ORDER BY Rank DESC")
	} else {
		searchString.WriteString("AND sprint in openSprints()")
	}

	return searchString.String()
}

func doRequest(request http.Request, data interface{}) error {
	resp, err := client.Do(&request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusNoContent {

		if err := json.Unmarshal(responseBody, data); err != nil {
			return err
		}
		return nil
	}

	return errors.New(string(responseBody))
}
