package cmd

import (
	"fmt"

	"github.com/fuskovic/gophercises/taskmanager/pkg/store"
	"github.com/spf13/cobra"
)

var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "display completed tasks",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := store.New()
		if err != nil {
			fmt.Printf("failed to init store : %s\n", err)
			return
		}
		defer s.Close()

		completedTasks, err := s.Completed()
		if err != nil {
			fmt.Printf("failed to list completed tasks : %s\n", err)
			return
		}

		fmt.Println("**********COMPLETED TASKS**********")

		for _, task := range completedTasks {
			fmt.Printf("%d - %s\n", task.ID, task.Description)
		}
	},
}

func init() {
	taskCmd.AddCommand(completedCmd)
}
