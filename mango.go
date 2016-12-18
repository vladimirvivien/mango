package mango

import (
	"reflect"

	"github.com/vladimirvivien/mango/gobuilder"
)

var (
	GoBuilder gobuilder.Builder
	tasks     []Task
)

func init() {
	tasks = []Task{}
	GoBuilder = gobuilder.New()
}

func Register(registrations ...interface{}) {
	if registrations == nil {
		return
	}

	for _, reg := range registrations {
		switch task := reg.(type) {
		case Task:
			tasks = append(tasks, task)
		case TaskFn:
			taskType := reflect.TypeOf(reg)

			tasks = append(tasks, Task{Do: task})
		default:
		}
	}
}

func Run(defaultTasks ...string) {

}
