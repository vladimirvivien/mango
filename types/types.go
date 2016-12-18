package mango

type TaskFn func()

type Task struct {
	Name string
	Desc string
	Do   TaskFn
}

type Builder interface {
	Build()
	Test()
}
