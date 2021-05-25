package mocks

import (
	"net/http"
)

// Thanks to https://www.thegreatcodeadventure.com/mocking-http-requests-in-golang/ for the approach

type MockRestClient struct{}
type MockCommandClient struct{}

var (
	DoFakeRequest  func(req *http.Request) (*http.Response, error)
	RunFakeCommand func(command string, args ...string) ([]byte, error)
)

func (m MockRestClient) Do(req *http.Request) (*http.Response, error) {
	return DoFakeRequest(req)
}

func (c MockCommandClient) Run(command string, args ...string) ([]byte, error) {
	return RunFakeCommand(command, args...)
}
