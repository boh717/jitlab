package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"

	"github.com/boh717/jitlab/pkg/command"
	"github.com/boh717/jitlab/pkg/git"
	"github.com/boh717/jitlab/pkg/gitlab"
	"github.com/boh717/jitlab/pkg/jira"
	"github.com/boh717/jitlab/pkg/question"
	"github.com/boh717/jitlab/pkg/rest"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile         string
	jiraService     jira.JiraService
	gitlabService   gitlab.GitlabService
	gitService      git.GitService
	questionService question.QuestionService
	rootCmd         = &cobra.Command{
		Use:     "jitlab",
		Short:   "Jitlab integrates Jira and GitLab for a faster development workflow",
		Long:    ``,
		Version: "0.1",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $HOME/.jitlab.json)")

	rootCmd.AddCommand(Config())
	rootCmd.AddCommand(InitRepo())
	rootCmd.AddCommand(NewTicket())
	rootCmd.AddCommand(Commits())
	rootCmd.AddCommand(MergeRequest())
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {

		home, err := homedir.Dir()
		cobra.CheckErr(err)

		viper.SetConfigFile(path.Join(home, ".jitlab.json"))
	}

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Fatalf("Config file %s not found", cfgFile)
	}

	jiraUrl := viper.GetString("jira.baseurl")
	validatedJiraBaseUrl, err := url.Parse(jiraUrl)
	if err != nil {
		log.Fatalf("Jira base URL %s is not a valid", jiraUrl)
	}
	gitlabUrl := viper.GetString("gitlab.baseurl")
	validatedGitlabBaseUrl, err := url.Parse(gitlabUrl)
	if err != nil {
		log.Fatalf("Gitlab base URL %s is not a valid", gitlabUrl)
	}
	jiraToken := viper.GetString("jira.token")
	jiraUsername := viper.GetString("jira.username")

	gitlabToken := viper.GetString("gitlab.token")
	gitlabGroup := viper.GetString("gitlab.groupid")

	branchPrefix := viper.GetString("branchPrefix")
	branchSuffix := viper.GetString("branchSuffix")
	keyCommitSeparator := viper.GetString("keyCommitSeparator")
	branchRegex := regexp.MustCompile(fmt.Sprintf("(%s)(\\w{1,6}-\\d{1,5})-(.*)(%s)", branchPrefix, branchSuffix))

	client := rest.RestClientImpl{Client: http.DefaultClient}
	commandClient := command.CommandClientImpl{}
	jiraService = jira.JiraServiceImpl{Client: client, BaseURL: validatedJiraBaseUrl.String(), Token: jiraToken, Username: jiraUsername}
	gitlabService = gitlab.GitlabServiceImpl{Client: client, BaseURL: validatedGitlabBaseUrl.String(), Token: gitlabToken, Group: gitlabGroup}
	gitService = git.GitServiceImpl{CommandClient: commandClient, BranchPrefix: branchPrefix, BranchSuffix: branchSuffix, KeyCommitSeparator: keyCommitSeparator, BranchRegexp: branchRegex}
	questionService = question.QuestionServiceImpl{}
}
