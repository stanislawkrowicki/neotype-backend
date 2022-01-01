all: web-api words

web-api:
	go run cmd/web-api/main.go

words:
	go run cmd/words/main.go

users:
	go run cmd/users/main.go

results-publisher:
	go run cmd/results-publisher/main.go

results-consumer:
	go run cmd/results-consumer/main.go

leaderboards:
	go run cmd/leaderboards/main.go

#DEBUG

web-api-debug:
	go build -gcflags="all=-N -l" -o ./bin/web-api ./cmd/web-api/main.go
	dlv --listen=localhost:40000 --headless=true --api-version=2 --only-same-user=false exec ./bin/web-api

words-debug:
	go build -gcflags="all=-N -l" -o ./bin/words ./cmd/words/main.go
	dlv --listen=localhost:40000 --headless=true --api-version=2 --only-same-user=false exec ./bin/words

users-debug:
	go build -gcflags="all=-N -l" -o ./bin/users ./cmd/users/main.go
	dlv --listen=localhost:40000 --headless=true --api-version=2 --only-same-user=false exec ./bin/users

results-publisher-debug:
	go build -gcflags="all=-N -l" -o ./bin/results-publisher ./cmd/results-publisher/main.go
	dlv --listen=localhost:40000 --headless=true --api-version=2 --only-same-user=false exec ./bin/results-publisher

results-consumer-debug:
	go build -gcflags="all=-N -l" -o ./bin/results-consumer ./cmd/results-consumer/main.go
	dlv --listen=localhost:40000 --headless=true --api-version=2 --only-same-user=false exec ./bin/results-consumer

leaderboards-debug:
	go build -gcflags="all=-N -l" -o ./bin/leaderboards ./cmd/leaderboards/main.go
	dlv --listen=localhost:40000 --headless=true --api-version=2 --only-same-user=false exec ./bin/leaderboards

#TESTS
test:
	go test ./tests/*

#DEPLOYMENT
deploy:
	sh ./deploy.sh