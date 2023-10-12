build:
	go build

run: build
	./petinfoservice

generate-proto:
	cd petinfoproto && make
generate-sqlc:
	sqlc generate
generate: generate-proto generate-sqlc

build-dockerfile:
	podman build . --tag claytontii/petinfoservice:latest
