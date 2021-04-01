package cmd

import (
	"log"

	"github.com/boh717/jitlab/pkg/question"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Config() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Configure Jitlab on first run",
		Long:  `Run this command the first time you run Jitlab to configure board and columns`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Configuring jitlab...")
			boards, err := jiraClient.GetBoards()
			if err != nil {
				log.Fatalln(err)
			}
			chosenBoard, err := question.AskForBoard(boards)
			if err != nil {
				log.Fatalln(err)
			}

			columns, err := jiraClient.GetBoardColumns(chosenBoard)
			if err != nil {
				log.Fatalln(err)
			}

			chosenColumns, err := question.AskForColumns(columns)
			if err != nil {
				log.Fatalln(err)
			}

			viper.Set("board", chosenBoard)
			viper.Set("columns", chosenColumns)

			if err := viper.WriteConfigAs(viper.ConfigFileUsed()); err != nil {
				log.Fatalln(err)
			}
		},
	}

	return configCmd

}
