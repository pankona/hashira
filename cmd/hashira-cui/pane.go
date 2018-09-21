package main

import (
	"fmt"
	"io"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/pankona/hashira/service"
)

type Pane struct {
	name       string
	index      int // place of this pane
	left       *Pane
	right      *Pane
	place      service.Place
	tasks      map[string]*service.Task
	priorities []string // array of task IDs
}

func (p *Pane) Layout(g *gocui.Gui, fucusedIndex int, selectedTask *service.Task) error {
	maxX, maxY := g.Size()
	v, err := g.SetView(p.name, maxX/4*p.index, 1, maxX/4*p.index+maxX/4-1, maxY-1)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = p.name
	}

	v.Clear()

	log.Printf("Tasks (%s): %v", p.place.String(), p.tasks)
	log.Printf("Priorities (%s): %v", p.place, p.priorities)
	return renderTasks(v, p.tasks, p.priorities, fucusedIndex, selectedTask)
}

func (p *Pane) len() int {
	return len(p.tasks)
}

func renderTasks(w io.Writer, tasks map[string]*service.Task, priorities []string, focusedIndex int, selectedTask *service.Task) error {
	var itemNum int
	var err error

	// render tasks for this pane
	for _, p := range priorities {
		task, ok := tasks[p]
		if !ok {
			continue
		}

		prefix := ""
		if selectedTask != nil && task.Id == selectedTask.Id {
			prefix = "*"
			log.Printf("selectedTask = %v", task)
		}
		if itemNum == focusedIndex {
			_, err = fmt.Fprintf(w, "%s \033[3%d;%dm%s\033[0m\n", prefix, 7, 4, task.Name)
			log.Printf("focusedIndex = %d", itemNum)
		} else {
			_, err = fmt.Fprintf(w, "%s %s\n", prefix, task.Name)
		}
		if err != nil {
			return err
		}
		itemNum++
	}

	return nil
}
