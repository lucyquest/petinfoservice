build:
	go build

run: build
	./petinfoservice

generate-proto:
	cd petinfoproto && make
generate: generate-proto

build-dockerfile:
	podman build . --tag claytontii/petinfoservice:latest
