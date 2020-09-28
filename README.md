### Golang Twitter scraping
Interacting with the Twitter API using Golang

### Running the project locally
- Export the required env vars 
```
export CONSUMER_KEY=<YOUR_VALUE_HERE>
export CONSUMER_SECRET=<YOUR_VALUE_HERE>
export ACCESS_TOKEN=<YOUR_VALUE_HERE>
export ACCESS_TOKEN_SECRET=<YOUR_VALUE_HERE>
export INTERVAL=<YOUR_VALUE_HERE> (MINUTES)
export QUERY_PARAM=<YOUR_VALUE_HERE>
export SEARCH_TYPE=<YOUR_VALUE_HERE>
export AZURE_KEY=<YOUR_VALUE_HERE>
export AZURE_ENDPOINT=<YOUR_VALUE_HERE>
```

- QUERY_PARAM - This is the query you wish to search/stream with
- SEARCH_TYPE - Either `STREAM` or `SEARCH` 
- AZURE_KEY - API key for the sentiment analysis service
- AZURE_ENDPOINT - Endpoint for the sentiment analysis service