package cmd

import (
	"fmt"
	"strconv"

	"github.com/fuskovic/gophercises/taskmanager/pkg/store"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove a task from incompleted tasks",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := store.New()
		if err != nil {
			fmt.Printf("failed to init store : %s\n", err)
			return
		}
		defer s.Close()

		if len(args) > 1 {
			fmt.Printf("the remove command only accepts one argument but received %d", len(args))
			return
		}

		ID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("failed to convert ID string input to int : %s\n", err)
			return
		}

		if err := s.Remove(ID); err != nil {
			fmt.Printf("failed to remove task ID : %d : %s\n", ID, err)
			return
		}
		fmt.Printf("successfully removed task ID %d!\n", ID)

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
	taskCmd.AddCommand(rmCmd)
}
