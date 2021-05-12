package jira

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/boh717/jitlab/pkg/rest"
)

type JiraService interface {
	GetBoards() ([]Board, error)
	GetBoardColumns(board Board) ([]Column, error)
	GetIssues(flowType string, projectKey string, columns []string, currentUser bool) ([]Issue, error)
}

type JiraServiceImpl struct {
	Client   rest.RestClient
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
	url := j.BaseURL + uri
	headers := map[string]string{"Authorization": fmt.Sprintf("Basic %s", j.Token)}

	req, err := j.Client.CreateRequest(http.MethodGet, url, headers, nil)
	if err != nil {
		return nil, err
	}

	response, err := j.Client.DoRequest(req)
	if err != nil {
		return nil, err
	}

	board := new(boardBase)
	err = j.Client.ProcessResponse(response, board)
	if err != nil {
		return nil, err
	}

	return board.Values, nil
}

func (j JiraServiceImpl) GetBoardColumns(board Board) ([]Column, error) {
	uri := fmt.Sprintf("/rest/agile/1.0/board/%d/configuration", board.ID)
	url := j.BaseURL + uri
	headers := map[string]string{"Authorization": fmt.Sprintf("Basic %s", j.Token)}

	req, err := j.Client.CreateRequest(http.MethodGet, url, headers, nil)
	if err != nil {
		return nil, err
	}

	response, err := j.Client.DoRequest(req)
	if err != nil {
		return nil, err
	}

	boardConfig := new(boardConfig)
	err = j.Client.ProcessResponse(response, boardConfig)
	if err != nil {
		return nil, err
	}

	return boardConfig.ColumnConfig.Columns, nil
}

func (j JiraServiceImpl) GetIssues(flowType string, projectKey string, columns []string, currentUser bool) ([]Issue, error) {
	searchString := buildSearchString(flowType, projectKey, columns, currentUser)
	uri := fmt.Sprintf("/rest/api/3/search?jql=%s&fields=summary,assignee&maxResults=50", url.QueryEscape(searchString))
	url := j.BaseURL + uri
	headers := map[string]string{"Authorization": fmt.Sprintf("Basic %s", j.Token)}

	req, err := j.Client.CreateRequest(http.MethodGet, url, headers, nil)
	if err != nil {
		return nil, err
	}

	response, err := j.Client.DoRequest(req)
	if err != nil {
		return nil, err
	}

	issue := new(issueBase)
	err = j.Client.ProcessResponse(response, issue)
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
