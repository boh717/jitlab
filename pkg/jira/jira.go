package jira

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var client = http.DefaultClient

type JiraService interface {
	GetBoards() ([]Board, error)
	GetBoardColumns(board Board) ([]Column, error)
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

func (j JiraServiceImpl) GetBoards() ([]Board, error) {
	uri := "/rest/agile/1.0/board"

	req, _ := http.NewRequest("GET", j.BaseURL+uri, nil)
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

	req, _ := http.NewRequest("GET", j.BaseURL+uri, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", j.Token))

	boardConfig := new(boardConfig)
	err := doRequest(*req, boardConfig)
	if err != nil {
		return nil, err
	}

	return boardConfig.ColumnConfig.Columns, nil
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
