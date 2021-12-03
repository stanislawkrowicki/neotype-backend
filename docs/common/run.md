## To run services, you need _make_ and _go1.17_.

### Running services 
#### To run a service, you can use the Makefile (while in main directory), eg:
    make web-api

#### If you want to run a service some other way, every entrypoint is located in _cmd_ directory.
    go run ./cmd/web-api/main.go



    