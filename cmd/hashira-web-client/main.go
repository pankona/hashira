package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	var (
		flagUpload   bool
		flagDownload bool
		flagTest     bool
	)

	flag.BoolVar(&flagUpload, "upload", false, "upload tasks and priorities to hashira-web")
	flag.BoolVar(&flagTest, "test", false, "test the hashira-web works")
	flag.Parse()

	accesstoken := os.Getenv("HASHIRA_ACCESS_TOKEN")
	if accesstoken == "" {
		log.Printf("Please specify environment variable HASHIRA_ACCESS_TOKEN. Abort.")
		os.Exit(1)
	}

	switch {
	case flagUpload:
		upload(accesstoken)
	case flagDownload:
		// TODO: implement
	case flagTest:
		fallthrough
	default:
		testAccessToken(accesstoken)
	}
}

const daemonPort = 50057

type Task struct {
	ID        string
	Name      string
	Place     string
	IsDeleted bool
}

type Priority map[string][]string
