package cmd

import (
	"fmt"

	"github.com/fuskovic/gophercises/taskmanager/pkg/store"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "display incomplete tasks",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := store.New()
		if err != nil {
			fmt.Printf("failed to init store : %s\n", err)
			return
		}
		defer s.Close()

		incompleteTasks, err := s.List()
		if err != nil {
			fmt.Printf("failed to list incomplete tasks : %s\n", err)
			return
		}

		fmt.Println("**********INCOMPLETE TASKS**********")

		for _, task := range incompleteTasks {
			fmt.Printf("%d - %s\n", task.ID, task.Description)
		}
	},
}

func init() {
	taskCmd.AddCommand(listCmd)
}
