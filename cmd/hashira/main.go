package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pankona/hashira/hashira/client"
)

func main() {
	ctx := context.Background()
	err := client.Create(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
