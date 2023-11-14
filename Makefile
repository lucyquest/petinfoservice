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

test: 
	$(info petinfoservice-makefile: Spawning postgres container)
	podman run --name petinfoservice-postgres --rm -e POSTGRES_PASSWORD=bestpassword -p 5432:5432 --health-cmd="pg_isready" postgres:16.0-alpine3.17 &

	until podman healthcheck run petinfoservice-postgres > /dev/null; do echo "petinfoservice-makefile: waiting on postgres" && sleep 0.1; done

	$(info petinfoservice-makefile: Starting Go Test)
	-POSTGRES_USER=postgres POSTGRES_PASS=bestpassword POSTGRES_HOST=localhost go test -race -v ./...

	$(info petinfoservice-makefile: Killing postgres)
	podman kill petinfoservice-postgres
