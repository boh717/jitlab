package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/boh717/jitlab/pkg/jira"
	"github.com/spf13/viper"
)

type GitService interface {
	CreateBranch(issue jira.Issue) error
}

type GitServiceImpl struct{}

func (g GitServiceImpl) CreateBranch(issue jira.Issue) error {
	replacer :=
		strings.NewReplacer(" ", "-", "~", "", "^", "", ":", "", "?", "", "*", "", "[", "", "]", "", "{", "", "}", "", "\\", "")

	summary := strings.ToLower(replacer.Replace(issue.Fields.Summary))
	issueKey := issue.Key
	prefix := viper.GetString("branchPrefix")
	suffix := viper.GetString("branchSuffix")

	branchName := fmt.Sprintf("%s%s-%s%s", prefix, issueKey, summary, suffix)

	command := exec.Command("git", "switch", "-c", branchName)
	return command.Run()

}
