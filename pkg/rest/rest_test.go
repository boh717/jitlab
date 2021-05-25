package rest_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/boh717/jitlab/pkg/mocks"
	"github.com/boh717/jitlab/pkg/rest"
	"github.com/google/go-cmp/cmp"
)

type person struct {
	FirstName string
	LastName  string
	Age       int
}

var successfulResponse = func(*http.Request) (*http.Response, error) {
	var successfulJson = `{"firstName":"Mario","lastName":"Rossi","age":25}`
	var successfulBody = ioutil.NopCloser(bytes.NewReader([]byte(successfulJson)))
	return &http.Response{
		StatusCode: 200,
		Body:       successfulBody,
	}, nil
}

var errorResponse = func(*http.Request) (*http.Response, error) {
	return nil, errors.New("Error from web server!")
}

var payload = "my payload"

func TestCreateRequest(t *testing.T) {
	tests := map[string]struct {
		method  string
		headers map[string]string
		payload io.Reader
	}{
		"GET without headers":              {method: "GET", headers: nil, payload: nil},
		"POST with headers and no payload": {method: "POST", headers: map[string]string{"MyToken": "secret-token"}, payload: nil},
		"POST with headers and payload":    {method: "GET", headers: map[string]string{"MyToken": "secret-token", "Content-Type": "application/json"}, payload: strings.NewReader(payload)},
	}
	mockHttpClient := mocks.MockRestClient{}
	restClient := rest.RestClientImpl{Client: mockHttpClient}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := restClient.CreateRequest(tc.method, "www.example.com", tc.headers, tc.payload)
			if err != nil {
				t.Fatalf("Got unexpected error: %v", err)
			}

			if tc.method != got.Method {
				t.Errorf("Passed method '%s' was not equal to '%s'", tc.method, got.Method)
			}
			for name, value := range tc.headers {
				setHeader := got.Header.Get(name)
				if setHeader != value {
					t.Errorf("Passed header '%s:%s' was not equal to '%s'", name, value, setHeader)
				}
			}

			if tc.payload == nil {

				if got.Body != nil {
					t.Error("Passed body was not empty!")
				}

			} else {

				gotPayload, _ := io.ReadAll(got.Body)
				gotPayloadString := string(gotPayload)

				if gotPayloadString != payload {
					t.Errorf("Passed payload '%s' was not equal to '%s'", payload, gotPayloadString)
				}
			}

		})
	}
}

func TestWrongCreateRequest(t *testing.T) {
	tests := map[string]struct {
		method string
		url    string
	}{
		"Bad method": {method: "Bad method", url: "http://www.example.com"},
		"Bad URL":    {method: "POST", url: " http://www.example.com"},
	}
	mockHttpClient := mocks.MockRestClient{}
	restClient := rest.RestClientImpl{Client: mockHttpClient}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := restClient.CreateRequest(tc.method, tc.url, nil, nil)
			if err == nil {
				t.Errorf("Got unexpected response: %+v", got)
			}

		})
	}
}

func TestDoRequest(t *testing.T) {

	tests := map[string]struct {
		externalOutput  func(*http.Request) (*http.Response, error)
		expectedSuccess bool
	}{
		"Successful request": {externalOutput: successfulResponse, expectedSuccess: true},
		"Wrong request":      {externalOutput: errorResponse, expectedSuccess: false},
	}
	mockHttpClient := mocks.MockRestClient{}
	restClient := rest.RestClientImpl{Client: mockHttpClient}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mocks.DoFakeRequest = tc.externalOutput
			request, _ := http.NewRequest("GET", "http://www.example.com", nil)

			resp, err := restClient.DoRequest(request)

			if tc.expectedSuccess {
				if err != nil {
					t.Errorf("Got error '%v', but wanted a success", err)
				}
			} else {
				if err == nil {
					t.Errorf("Got response '%+v', but wanted an error", resp)
				}
			}

		})
	}
}

func TestProcessResponse(t *testing.T) {

	tests := map[string]struct {
		statusCode      int
		body            string
		want            *person
		expectedSuccess bool
	}{
		"Successful path":    {statusCode: 200, body: `{"firstName":"Mario","lastName":"Rossi","age":25}`, want: &person{FirstName: "Mario", LastName: "Rossi", Age: 25}, expectedSuccess: true},
		"Malformed json":     {statusCode: 200, body: `{"firstName""Mario","lastName":"Rossi","age":25}`, want: nil, expectedSuccess: false},
		"Resource not found": {statusCode: 404, body: "Resource not found", want: nil, expectedSuccess: false},
	}
	mockHttpClient := mocks.MockRestClient{}
	restClient := rest.RestClientImpl{Client: mockHttpClient}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			person := new(person)
			response := &http.Response{StatusCode: tc.statusCode, Body: ioutil.NopCloser(bytes.NewReader([]byte(tc.body)))}
			err := restClient.ProcessResponse(response, person)

			if tc.expectedSuccess {
				if err != nil {
					t.Errorf("Got unexpected error %v", err)
				}

				if !cmp.Equal(person, tc.want) {
					t.Errorf("Result person '%+v' is different from expected one '%+v'", person, tc.want)
				}
			} else {
				if err == nil {
					t.Errorf("Got response '%+v', expected error", person)
				}
			}

		})
	}
}
