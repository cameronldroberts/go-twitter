test: 
	go test ./...
test-cover:
	go test ./... -coverprofile=c.out && go tool cover -html=c.out