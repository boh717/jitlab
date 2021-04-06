package gitlab

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var client = http.DefaultClient

type GitlabService interface {
	SearchProject(search string) ([]Repository, error)
}

type GitlabServiceImpl struct {
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

func (g GitlabServiceImpl) SearchProject(search string) ([]Repository, error) {
	uri := fmt.Sprintf("/groups/%s/search?scope=projects&search=%s", g.Group, search)

	req, _ := http.NewRequest("GET", g.BaseURL+uri, nil)
	req.Header.Set("PRIVATE-TOKEN", g.Token)

	var repositories []Repository
	err := doRequest(*req, &repositories)
	if err != nil {
		return nil, err
	}

	return repositories, nil
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
