build:
	CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/api/main.go

start_dev: 
	export GO_ENV=development && go run ./cmd/api/main.go