package mocks

import (
	"net/http"
)

// Thanks to https://www.thegreatcodeadventure.com/mocking-http-requests-in-golang/

type MockClient struct{}

var (
	DoFakeRequest func(req *http.Request) (*http.Response, error)
)

func (m MockClient) Do(req *http.Request) (*http.Response, error) {
	return DoFakeRequest(req)
}
