package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/boh717/jitlab/pkg/jira"
)

type GitService interface {
	GetCurrentBranch() (string, error)
	CreateBranch(issue jira.Issue) error
	CreateTitleFromBranch(branch string) (string, error)
	Commit(message string) error
	Push(branch string) error
}

type GitServiceImpl struct {
	BranchPrefix       string
	BranchSuffix       string
	KeyCommitSeparator string
	BranchRegexp       *regexp.Regexp
}

func (g GitServiceImpl) GetCurrentBranch() (string, error) {
	var out bytes.Buffer
	replacer := strings.NewReplacer(" ", "", "\n", "")

	command := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	command.Stdout = &out
	err := command.Run()
	if err != nil {
		return "", err
	}

	return replacer.Replace(out.String()), nil
}

func (g GitServiceImpl) CreateBranch(issue jira.Issue) error {
	replacer :=
		strings.NewReplacer(" ", "-", "~", "", "^", "", ":", "", "?", "", "*", "", "[", "", "]", "", "{", "", "}", "", "\\", "")

	summary := strings.ToLower(replacer.Replace(issue.Fields.Summary))
	issueKey := issue.Key

	branchName := fmt.Sprintf("%s%s-%s%s", g.BranchPrefix, issueKey, summary, g.BranchSuffix)

	command := exec.Command("git", "switch", "-c", branchName)
	return command.Run()

}

func (g GitServiceImpl) CreateTitleFromBranch(branch string) (string, error) {
	key := getIssueKeyFromBranch(branch, g.BranchRegexp)
	title := getIssueTitleFromBranch(branch, g.BranchRegexp)

	prettyTitle := strings.ReplaceAll(title, "-", " ")

	switch {
	case key != "" && prettyTitle != "":
		return fmt.Sprintf("%s: %s", key, prettyTitle), nil

	case key == "" && prettyTitle != "":
		return prettyTitle, nil

	default:
		return strings.ReplaceAll(branch, "-", " "), nil
	}

}

func (g GitServiceImpl) Commit(message string) error {
	var commitMessage string
	branch, err := g.GetCurrentBranch()
	if err != nil {
		return err
	}

	key := getIssueKeyFromBranch(branch, g.BranchRegexp)

	commitMessage = message
	if key != "" {
		commitMessage = fmt.Sprintf("%s%s%s", key, g.KeyCommitSeparator, commitMessage)
	}

	command := exec.Command("git", "commit", "-m", commitMessage)
	return command.Run()
}

func (g GitServiceImpl) Push(branch string) error {
	command := exec.Command("git", "push", "--set-upstream", "origin", branch)
	return command.Run()
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
