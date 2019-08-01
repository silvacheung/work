package taskpool

import (
	"strings"
)

type Exec func()

type Task struct {
	name string
	exec Exec
}

func New(name string, exec Exec) *Task {
	if len(strings.Trim(name, " ")) == 0 {
		name = DefaultTask
	}
	return &Task{name: name, exec: exec}
}

func (t *Task) Exec() {
	if t.exec != nil {
		t.exec()
	}
}
