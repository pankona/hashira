package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func testAccessToken(accesstoken string) {
	req, err := http.NewRequest(http.MethodGet, "https://asia-northeast1-hashira-web.cloudfunctions.net/test-access-token", nil)
	if err != nil {
		log.Printf("failed to create new request: %v", err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", accesstoken))

	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		log.Printf("request failed: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.Printf("This accesstoken is not valid. Please check HASHIRA_ACCESS_TOKEN is correct and try again [%d].", resp.StatusCode)
		os.Exit(1)
	}

	log.Printf("hashira-web works!")
}
