package processing

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.1/textanalytics"
	"github.com/Azure/go-autorest/autorest/to"

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
			SentimentAnalysis(tweets)
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

	ch := make(chan os.Signal,1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()
}

func SentimentAnalysis(tweetBatch []TweetData) {

	textAnalyticsClient := auth.GetTextAnalyticsClient()
	ctx := context.Background()
	inputDocuments := returnInputDocuments(tweetBatch)

	batchInput := textanalytics.MultiLanguageBatchInput{Documents: &inputDocuments}
	result, _ := textAnalyticsClient.Sentiment(ctx, to.BoolPtr(false), &batchInput)
	var batchResult textanalytics.SentimentBatchResult
	jsonString, _ := json.Marshal(result)
	_ = json.Unmarshal(jsonString, &batchResult)

	// Printing sentiment results
	for _, document := range *batchResult.Documents {
		logr.Printf("Document ID: %s ", *document.ID)
		logr.Printf("Sentiment Score: %f\n", *document.Score)
	}

	// Printing document errors
	for _, err := range *batchResult.Errors {
		logr.Printf("Document ID: %s Message : %s\n", *err.ID, *err.Message)
	}
}

func returnInputDocuments(batch []TweetData) []textanalytics.MultiLanguageInput {
	tweetsToProcess := []textanalytics.MultiLanguageInput{}
	for i, v := range batch {
		test := textanalytics.MultiLanguageInput{
			Language: to.StringPtr("en"),
			ID:       to.StringPtr(strconv.Itoa(i)),
			Text:     to.StringPtr(v.Body),
		}
		logr.Println(i, " : ", v.Body)
		tweetsToProcess = append(tweetsToProcess, test)
	}

	return tweetsToProcess
}
