package git

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/boh717/jitlab/pkg/command"
	"github.com/boh717/jitlab/pkg/jira"
)

type GitService interface {
	GetCurrentBranch() (string, error)
	CreateBranch(issue jira.Issue) (string, error)
	CreateTitleFromBranch(branch string) (string, error)
	Commit(branch string, message string) (string, error)
	Push(branch string) (string, error)
}

type GitServiceImpl struct {
	CommandClient      command.CommandClient
	BranchPrefix       string
	BranchSuffix       string
	KeyCommitSeparator string
	BranchRegexp       *regexp.Regexp
}

func (g GitServiceImpl) GetCurrentBranch() (string, error) {
	replacer := strings.NewReplacer(" ", "", "\n", "")

	out, err := g.CommandClient.Run("git", "branch", "--show-current")
	if err != nil {
		return "", errors.New(fmt.Sprint(err) + ": " + string(out))
	}

	return replacer.Replace(string(out)), nil
}

func (g GitServiceImpl) CreateBranch(issue jira.Issue) (string, error) {
	replacer :=
		strings.NewReplacer(" ", "-", "~", "", "^", "", ":", "", "?", "", "*", "", "[", "", "]", "", "{", "", "}", "", "\\", "")

	summary := strings.ToLower(replacer.Replace(issue.Fields.Summary))
	issueKey := issue.Key

	branchName := fmt.Sprintf("%s%s-%s%s", g.BranchPrefix, issueKey, summary, g.BranchSuffix)

	out, err := g.CommandClient.Run("git", "switch", "-c", branchName)
	if err != nil {
		return "", errors.New(fmt.Sprint(err) + ": " + string(out))
	}

	return branchName, nil

}

func (g GitServiceImpl) CreateTitleFromBranch(branch string) (string, error) {
	key := getIssueKeyFromBranch(branch, g.BranchRegexp)
	title := getIssueTitleFromBranch(branch, g.BranchRegexp)

	prettyTitle := strings.ReplaceAll(title, "-", " ")

	switch {
	case key != "" && prettyTitle != "":
		return fmt.Sprintf("%s%s %s", key, g.KeyCommitSeparator, prettyTitle), nil

	case key == "" && prettyTitle != "":
		return prettyTitle, nil

	default:
		return strings.ReplaceAll(branch, "-", " "), nil
	}

}

func (g GitServiceImpl) Commit(branch string, message string) (string, error) {
	var commitMessage string

	key := getIssueKeyFromBranch(branch, g.BranchRegexp)

	commitMessage = message
	if key != "" {
		commitMessage = fmt.Sprintf("%s%s %s", key, g.KeyCommitSeparator, commitMessage)
	}

	out, err := g.CommandClient.Run("git", "commit", "-m", commitMessage)
	if err != nil {
		return "", errors.New(fmt.Sprint(err) + ": " + string(out))
	}

	return commitMessage, nil

}

func (g GitServiceImpl) Push(branch string) (string, error) {
	out, err := g.CommandClient.Run("git", "push", "--set-upstream", "origin", branch)
	if err != nil {
		return "", errors.New(fmt.Sprint(err) + ": " + string(out))
	}

	return string(out), nil

}

func getIssueKeyFromBranch(branch string, r *regexp.Regexp) string {
	matches := r.FindStringSubmatch(branch)
	if matches != nil {
		return matches[2]
	}

	return ""
}

func getIssueTitleFromBranch(branch string, r *regexp.Regexp) string {
	matches := r.FindStringSubmatch(branch)
	if matches != nil {
		return matches[3]
	}

	return ""
}
