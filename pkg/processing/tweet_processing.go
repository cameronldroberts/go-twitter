package processing

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	auth "github.com/cameronldroberts/go-twitter/pkg/auth"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	logr "github.com/sirupsen/logrus"
)

type TweetData struct {
	Body         string
	LikeCount    int
	RetweetCount int
}

func SearchTweet(client *twitter.Client) {
	queryParam := os.Getenv("QUERY_PARAM")
	search, resp, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: queryParam,
	})

	if err != nil {
		log.Print(err)
	}

	log.Printf("%+v\n", resp.Body)
	log.Printf("%+v\n", search)
}

func TweetStream() {
	creds := auth.GetCreds()
	queryParam := os.Getenv("QUERY_PARAM")
	logr.Print("Getting tweet stream for : ", queryParam)
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	var tweets []TweetData

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		tweetdata := TweetData{
			Body:         tweet.Text,
			LikeCount:    tweet.FavoriteCount,
			RetweetCount: tweet.RetweetCount,
		}

		tweets = append(tweets, tweetdata)
		if len(tweets) == 10 {
			// Call sentiment API in batches
			logr.Print(tweets)
			for _, tweets := range tweets {
				logr.Print(tweets)
			}
			tweets = nil
		}
	}

	filterParams := &twitter.StreamFilterParams{
		Track:         []string{queryParam},
		StallWarnings: twitter.Bool(true),
		Language:      []string{"en"},
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	go demux.HandleChan(stream.Messages)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()
}
