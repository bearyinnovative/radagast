package task

import "github.com/bearyinnovative/radagast/monitor_stale_issues"

var AvailableTasks = map[string]TaskExecuteFn{
	monitor_stale_issues.TaskName: monitor_stale_issues.Execute,
}
