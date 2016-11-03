package task

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/bearyinnovative/radagast/config"
)

var (
	ErrNoTasksDefined = errors.New("no tasks defined")
)

// Execute enabled tasks. Blocks until all tasks is finished.
func Execute(c context.Context) error {
	config := config.FromContext(c)

	tasks, err := getTasks(config)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}

	for _, taskName := range tasks {
		if err := spawnTask(c, taskName, wg); err != nil {
			return err
		}
	}

	wg.Wait()

	return nil
}

func spawnTask(c context.Context, taskName string, wg *sync.WaitGroup) error {
	taskExecutor, present := AvailableTasks[taskName]
	if !present {
		return fmt.Errorf("unknown task: %s", taskName)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := taskExecutor(c); err != nil {
			log.Printf("execute task %s failed: %+v", taskName, err)
		}
	}()

	return nil
}

func getTasks(c map[string]interface{}) (tasks []string, err error) {
	itasks, ok := c["tasks"].([]interface{})
	if !ok {
		err = ErrNoTasksDefined
		return
	}

	for _, itask := range itasks {
		task := itask.(string)
		tasks = append(tasks, task)
	}

	return
}
