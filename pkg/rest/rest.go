package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type RestClient interface {
	DoRequest(req *http.Request) (*http.Response, error)
	ProcessRequest(resp *http.Response, data interface{}) error
}

type RestClientImpl struct {
	Client http.Client
}

func (r RestClientImpl) DoRequest(request *http.Request) (*http.Response, error) {
	return r.Client.Do(request)
}

func (r RestClientImpl) ProcessRequest(resp *http.Response, data interface{}) error {
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
