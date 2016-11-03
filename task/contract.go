package task

import "context"

// TaskExecuteFn executes a task.
type TaskExecuteFn func(context.Context) error
