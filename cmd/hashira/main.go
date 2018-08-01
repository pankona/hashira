package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/pankona/hashira/hashira/client"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	c := client.Client{
		Address: "localhost:50056",
	}
	ctx := context.Background()

	newCmd := kingpin.Command("new", "add new task with specified task name")
	var (
		name = newCmd.Arg("name", "name of task which is newly added").Required().String()
		_    = newCmd.Action(func(pc *kingpin.ParseContext) error {
			err := c.Create(ctx, *name)
			if err != nil {
				return errors.New("failed to create a new task: " + err.Error())
			}
			return nil
		})
	)

	kingpin.Parse()

	tasks, err := c.Retrieve(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, v := range tasks {
		id, err := strconv.Atoi(v.Id)
		if err != nil {
			continue
		}
		fmt.Printf("[%04d]\t%s\t%v\t%v\n", id, v.Name, v.Place, v.IsDeleted)
	}
}
