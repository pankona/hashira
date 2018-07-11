package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pankona/hashira/hashira/client"
)

func main() {
	c := client.Client{
		Address: "localhost:50056",
	}
	ctx := context.Background()
	err := c.Create(ctx, "test")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
