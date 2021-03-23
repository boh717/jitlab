package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func InitRepo() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Configure your repository",
		Long:  `Run this command in every git repo you want to use Jitlab`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Init repo...")
		},
	}

	return initCmd
}
