package store

type (
	// Task represents a to-do item that has or has not been completed.
	Task struct {
		ID          int
		Description string
		IsCompleted bool
	}
	// Tasks contains a list of tasks.
	Tasks []Task
)
