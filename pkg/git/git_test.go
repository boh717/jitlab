package git_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/boh717/jitlab/pkg/git"
	"github.com/boh717/jitlab/pkg/jira"
	"github.com/boh717/jitlab/pkg/mocks"
)

var branchName = "your-branch-name"

func TestGetCurrentBranch(t *testing.T) {
	tests := map[string]struct {
		command        func(command string, args ...string) ([]byte, error)
		expectedBranch string
	}{
		"Return current branch": {func(command string, args ...string) ([]byte, error) { return []byte(branchName), nil }, branchName},
		"Return branch correctly stripped": {func(command string, args ...string) ([]byte, error) {
			return []byte(fmt.Sprintf("\n%s \n", branchName)), nil
		}, branchName},
		"Return error": {func(command string, args ...string) ([]byte, error) { return nil, errors.New("Fatal!") }, ""},
	}
	mockCommandClient := mocks.MockCommandClient{}
	gitClient := git.GitServiceImpl{CommandClient: mockCommandClient}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mocks.RunFakeCommand = tc.command
			result, err := gitClient.GetCurrentBranch()

			if tc.expectedBranch == "" && err == nil {
				t.Errorf("Got no branch nor error. Something unexpected happened!")
			}

			if tc.expectedBranch != result {
				t.Errorf("Wanted branch '%s'. Got branch '%s' instead", tc.expectedBranch, result)
			}

		})
	}
}

func TestCreateBranch(t *testing.T) {
	tests := map[string]struct {
		command        func(command string, args ...string) ([]byte, error)
		issue          jira.Issue
		expectedBranch string
	}{
		"Return created branch": {func(command string, args ...string) ([]byte, error) {
			return []byte("Success"), nil
		}, initJiraIssue("JT-01", "Complete this task"), "prefix/JT-01-complete-this-task-suffix"},
		"Return created branch correctly stripped": {func(command string, args ...string) ([]byte, error) {
			return []byte("Success"), nil
		}, initJiraIssue("JT-01", "[POC] Complete this task maybe?"), "prefix/JT-01-poc-complete-this-task-maybe-suffix"},
		"Return error": {func(command string, args ...string) ([]byte, error) {
			return nil, errors.New("Failed creating branch")
		}, initJiraIssue("JT-01", "[POC] Complete this task maybe?"), ""},
	}
	mockCommandClient := mocks.MockCommandClient{}
	gitClient := git.GitServiceImpl{
		CommandClient:      mockCommandClient,
		BranchPrefix:       "prefix/",
		BranchSuffix:       "-suffix",
		KeyCommitSeparator: ":",
		BranchRegexp:       regexp.MustCompile(fmt.Sprintf("(%s)(\\w{1,6}-\\d{1,4})-(.*)(%s)", "prefix", "suffix")),
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mocks.RunFakeCommand = tc.command
			result, err := gitClient.CreateBranch(tc.issue)

			if tc.expectedBranch == "" && err == nil {
				t.Errorf("Got no branch nor error. Something unexpected happened!")
			}

			if tc.expectedBranch != result {
				t.Errorf("Wanted branch '%s'. Got branch '%s' instead", tc.expectedBranch, result)
			}

		})
	}
}

func TestCreateTitleFromBranch(t *testing.T) {
	tests := map[string]struct {
		prefix             string
		suffix             string
		keyCommitSeparator string
		branch             string
		expectedTitle      string
	}{
		"Return title from complete branch":                  {prefix: "prefix/", suffix: "-suffix", keyCommitSeparator: ":", branch: "prefix/JT-01-complete-this-task-suffix", expectedTitle: "JT-01: complete this task"},
		"Return title from branch without prefix":            {suffix: "-suffix", keyCommitSeparator: ":", branch: "JT-01-complete-this-task-suffix", expectedTitle: "JT-01: complete this task"},
		"Return title from branch without suffix":            {prefix: "prefix/", keyCommitSeparator: ":", branch: "prefix/JT-01-complete-this-task", expectedTitle: "JT-01: complete this task"},
		"Return title from branch without prefix nor suffix": {keyCommitSeparator: "#", branch: "JT-01-complete-this-task", expectedTitle: "JT-01# complete this task"},
		"Return title from branch with unexpected prefix":    {keyCommitSeparator: ":", branch: "prefix/JT-01-complete-this-task", expectedTitle: "JT-01: complete this task"},
		"Return title from branch without key":               {keyCommitSeparator: ":", branch: "complete-this-task", expectedTitle: "complete this task"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gitClient := git.GitServiceImpl{
				BranchPrefix:       tc.prefix,
				BranchSuffix:       tc.suffix,
				KeyCommitSeparator: tc.keyCommitSeparator,
				BranchRegexp:       regexp.MustCompile(fmt.Sprintf("(%s)(\\w{1,6}-\\d{1,5})-(.*)(%s)", tc.prefix, tc.suffix)),
			}
			result, _ := gitClient.CreateTitleFromBranch(tc.branch)

			if tc.expectedTitle != result {
				t.Errorf("Got title '%s', but wanted '%s'", result, tc.expectedTitle)
			}

		})
	}
}

func TestCommit(t *testing.T) {
	tests := map[string]struct {
		command           func(command string, args ...string) ([]byte, error)
		branch            string
		message           string
		expectedCommitMsg string
	}{
		"Return commit message": {func(command string, args ...string) ([]byte, error) {
			return []byte("Success"), nil
		}, "prefix/JT-01-complete-this-task", "Add feature X", "JT-01: Add feature X"},
		"Return error": {func(command string, args ...string) ([]byte, error) {
			return nil, errors.New("Fatal!")
		}, "prefix/JT-01-complete-this-task", "Add feature X", ""},
	}
	mockCommandClient := mocks.MockCommandClient{}
	gitClient := git.GitServiceImpl{
		CommandClient:      mockCommandClient,
		BranchPrefix:       "prefix/",
		KeyCommitSeparator: ":",
		BranchRegexp:       regexp.MustCompile(fmt.Sprintf("(%s)(\\w{1,6}-\\d{1,4})-(.*)(%s)", "prefix/", "")),
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mocks.RunFakeCommand = tc.command
			result, err := gitClient.Commit(tc.branch, tc.message)

			if tc.expectedCommitMsg == "" && err == nil {
				t.Errorf("Got no message nor error. Something unexpected happened!")
			}

			if result != tc.expectedCommitMsg {
				t.Errorf("Wanted commit message '%s'. Got message '%s' instead", tc.expectedCommitMsg, result)
			}

		})
	}
}

func TestPush(t *testing.T) {
	tests := map[string]struct {
		command         func(command string, args ...string) ([]byte, error)
		branch          string
		expectedPushMsg string
	}{
		"Return push message": {func(command string, args ...string) ([]byte, error) { return []byte("Success"), nil }, "prefix/JT-01-complete-this-task", "Success"},
		"Return push error":   {func(command string, args ...string) ([]byte, error) { return nil, errors.New("Fatal!") }, "prefix/JT-01-complete-this-task", ""},
	}
	mockCommandClient := mocks.MockCommandClient{}
	gitClient := git.GitServiceImpl{CommandClient: mockCommandClient}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mocks.RunFakeCommand = tc.command
			result, err := gitClient.Push(tc.branch)

			if tc.expectedPushMsg == "" && err == nil {
				t.Errorf("Got no push message nor error. Something unexpected happened!")
			}

			if result != tc.expectedPushMsg {
				t.Errorf("Wanted push message '%s'. Got message '%s' instead", tc.expectedPushMsg, result)
			}

		})
	}
}

func initJiraIssue(key string, summary string) jira.Issue {
	issue := jira.Issue{}
	issue.ID = "id"
	issue.Key = key
	issue.Fields.Summary = summary

	return issue
}
