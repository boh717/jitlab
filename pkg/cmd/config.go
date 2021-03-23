package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Config() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Configure Jitlab on first run",
		Long:  `Run this command the first time you run Jitlab to configure board and input/output columns`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Configuration...")
		},
	}

	return configCmd

}
