build:
	go build

run: build
	./petinfoservice

precommit:
	golangci-lint run
	gosec ./...
	go mod tidy

generate-proto:
	cd petinfoproto && make
generate-sqlc:
	sqlc generate
generate: generate-proto generate-sqlc

build-dockerfile:
	podman build . --tag claytontii/petinfoservice:latest

run-postgres:
	podman run --name petinfoservice-postgres --rm -e POSTGRES_PASSWORD=bestpassword -it -p 5432:5432 postgres:16.0-alpine3.17
