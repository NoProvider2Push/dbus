build-prod:
	go mod download
					GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/np2p-gnu-linux-amd64 . &\
	CC=musl-gcc 			GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/np2p-alpine-linux-amd64 . &\
					GOOS=linux GOARCH=386 CGO_ENABLED=1 go build -o bin/np2p-gnu-linux-386 . &\
	CC=aarch64-linux-gnu-gcc 	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -o bin/np2p-gnu-linux-arm64 . &\
	CC=aarch64-linux-musl-gcc 	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -o bin/np2p-alpine-linux-arm64 .

archlinux-deps-prod: archlinux-deps-init archlinux-deps

archlinux-deps-init:
	sudo pacman -Syu --needed --noconfirm git base-devel
	mkdir -p ~/Downloads; \
		cd ~/Downloads; \
		git clone https://aur.archlinux.org/yay-bin.git; \
		cd yay-bin; \
		makepkg -si --noconfirm
archlinux-deps:
	#                           aarch64 gnu -      amd64 alpine - aarch64 alpine - 32bit gnu - 32 bit gnu - go
	yay --noconfirm --needed -S aarch64-linux-gnu-gcc musl aarch64-linux-musl lib32-glibc lib32-gcc-libs go
