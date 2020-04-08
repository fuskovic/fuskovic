package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var taskCmd = cobra.Command{
	Use:   "task",
	Short: "CLI task-manager",
	Run:   func(cmd *cobra.Command, args []string) { cmd.Help() },
}

// Execute runs task
func Execute() {
	if err := taskCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
