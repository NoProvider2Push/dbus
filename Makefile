build-local:
	go build -o np2p_dbus
build-prod:
	go mod download
					GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/np2p-gnu-linux-amd64 . &\
	CC=musl-gcc 			GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/np2p-alpine-linux-amd64 . &\
					GOOS=linux GOARCH=386 CGO_ENABLED=1 go build -o bin/np2p-gnu-linux-386 . &\
	CC=aarch64-linux-gnu-gcc 	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -o bin/np2p-gnu-linux-arm64 . &\
	CC=aarch64-linux-musl-gcc 	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -o bin/np2p-alpine-linux-arm64 .

test: build-local
	go test ./...

build-docker:
	mkdir bin
	chmod ugo+rwx bin
	docker run --rm -v `pwd`:/app ghcr.io/karmanyaahm/mega_go_arch_xcompiler:v0.2.1 build np2p
