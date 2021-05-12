package rest

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

type RestClient interface {
	CreateRequest(method string, url string, headers map[string]string, payload io.Reader) (*http.Request, error)
	DoRequest(req *http.Request) (*http.Response, error)
	ProcessResponse(resp *http.Response, data interface{}) error
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type RestClientImpl struct {
	Client httpClient
}

func (r RestClientImpl) CreateRequest(method string, url string, headers map[string]string, payload io.Reader) (*http.Request, error) {

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	for key, element := range headers {
		req.Header.Set(key, element)
	}

	return req, nil
}

func (r RestClientImpl) DoRequest(req *http.Request) (*http.Response, error) {

	return r.Client.Do(req)
}

func (r RestClientImpl) ProcessResponse(resp *http.Response, data interface{}) error {
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
