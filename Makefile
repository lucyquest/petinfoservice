build:
	go build
generate: generate-proto
generate-proto:
	cd petinfoproto && make
build-dockerfile:
	podman build . --tag claytontii/petinfoservice:latest
