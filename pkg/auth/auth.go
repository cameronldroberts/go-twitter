package auth

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.1/textanalytics"
	"github.com/Azure/go-autorest/autorest"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

func GetCreds() Credentials {
	creds := Credentials{
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
	}
	return creds
}

func GetClient(creds *Credentials) (*twitter.Client, error) {

	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// verifyParams := &twitter.AccountVerifyParams{
	// 	SkipStatus:   twitter.Bool(true),
	// 	IncludeEmail: twitter.Bool(true),
	// }

	// _, _, err := client.Accounts.VerifyCredentials(verifyParams)
	// if err != nil {
	// 	return nil, err
	// }

	return client, nil
}

func GetTextAnalyticsClient() textanalytics.BaseClient {
	key := os.Getenv("AZURE_KEY")
	endpoint := os.Getenv("AZURE_ENDPOINT")

	textAnalyticsClient := textanalytics.New(endpoint)
	textAnalyticsClient.Authorizer = autorest.NewCognitiveServicesAuthorizer(key)

	return textAnalyticsClient
}
