all: web-api words

web-api:
	go run cmd/web-api/main.go

words:
	go run cmd/words/main.go

#DEBUG

web-api-debug:
	go build -gcflags="all=-N -l" -o ./bin/web-api ./cmd/web-api/main.go
	dlv --listen=localhost:40000 --headless=true --api-version=2 --only-same-user=false exec ./bin/web-api

words-debug:
	go build -gcflags="all=-N -l" -o ./bin/words ./cmd/words/main.go
	dlv --listen=localhost:40000 --headless=true --api-version=2 --only-same-user=false exec ./bin/words