package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/cameronldroberts/go-twitter/pkg/auth"
	"github.com/cameronldroberts/go-twitter/pkg/processing"
	logr "github.com/sirupsen/logrus"
)

func main() {
	searchType := os.Getenv("SEARCH_TYPE")
	if searchType == "SEARCH" {
		logr.Println("About to start tweet search")
		var interval string = os.Getenv("INTERVAL")
		creds := auth.GetCreds()
		client, err := auth.GetClient(&creds)

		if err != nil {
			log.Println("Error getting Twitter Client")
			log.Println(err)
		}

		for {
			processing.SearchTweet(client)
			sleepInterval, err := strconv.Atoi(interval)
			if err != nil {
				fmt.Println("Error converting string", err)
			}

			time.Sleep(time.Duration(sleepInterval) * time.Second)
		}
	} else if searchType == "STREAM" {
		logr.Println("About to start tweet stream")
		processing.TweetStream()
	} else if searchType == "SENTI" {
		processing.TweetStream()
	}
}
