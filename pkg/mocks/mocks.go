package mocks

import (
	"net/http"
)

// Thanks to https://www.thegreatcodeadventure.com/mocking-http-requests-in-golang/ for the approach

type MockClient struct{}
type CmdClient struct{}

var (
	DoFakeRequest  func(req *http.Request) (*http.Response, error)
	RunFakeCommand func() error
)

func (m MockClient) Do(req *http.Request) (*http.Response, error) {
	return DoFakeRequest(req)
}

func (c CmdClient) Run() error {
	return RunFakeCommand()
}
