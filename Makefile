all: web-api words

web-api:
	go run cmd/web-api/main.go

words:
	go run cmd/words/main.go