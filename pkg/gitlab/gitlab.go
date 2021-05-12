package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boh717/jitlab/pkg/rest"
)

type GitlabService interface {
	SearchProject(search string) ([]Repository, error)
	CreateMergeRequest(projectId string, sourceBranch string, targetBranch string, title string, removeSourceBranch bool, squash bool) (MergeRequestResponse, error)
}

type GitlabServiceImpl struct {
	Client  rest.RestClient
	BaseURL string
	Token   string
	Group   string
}

type Repository struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Path        string `json:"path"`
}

type mrRequest struct {
	ID                 string `json:"id"`
	SourceBranch       string `json:"source_branch"`
	TargetBranch       string `json:"target_branch"`
	Title              string `json:"title"`
	RemoveSourceBranch bool   `json:"remove_source_branch"`
	Squash             bool   `json:"squash"`
}

type MergeRequestResponse struct {
	Url string `json:"web_url"`
}

func (g GitlabServiceImpl) SearchProject(search string) ([]Repository, error) {
	uri := fmt.Sprintf("/groups/%s/search?scope=projects&search=%s", g.Group, search)
	url := g.BaseURL + uri
	headers := map[string]string{"PRIVATE-TOKEN": g.Token}

	req, err := g.Client.CreateRequest(http.MethodGet, url, headers, nil)
	if err != nil {
		return nil, err
	}

	response, err := g.Client.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var repositories []Repository
	err = g.Client.ProcessResponse(response, repositories)
	if err != nil {
		return nil, err
	}

	return repositories, nil
}

func (g GitlabServiceImpl) CreateMergeRequest(projectId string, sourceBranch string, targetBranch string, title string, removeSourceBranch bool, squash bool) (MergeRequestResponse, error) {
	mrResponse := MergeRequestResponse{}
	uri := fmt.Sprintf("/projects/%s/merge_requests", projectId)
	url := g.BaseURL + uri
	headers := map[string]string{"PRIVATE-TOKEN": g.Token, "Content-Type": "application/json"}
	request := mrRequest{
		ID:                 projectId,
		SourceBranch:       sourceBranch,
		TargetBranch:       targetBranch,
		Title:              title,
		RemoveSourceBranch: removeSourceBranch,
		Squash:             squash}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return mrResponse, err
	}

	req, err := g.Client.CreateRequest(http.MethodPost, url, headers, bytes.NewBuffer(jsonRequest))
	if err != nil {
		return mrResponse, err
	}

	response, err := g.Client.DoRequest(req)
	if err != nil {
		return mrResponse, err
	}

	err = g.Client.ProcessResponse(response, mrResponse)
	if err != nil {
		return mrResponse, err
	}

	return mrResponse, nil
}
