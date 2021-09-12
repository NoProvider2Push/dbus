build-prod:
					GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/np2p-gnu-linux-amd64 .
	CC=musl-gcc 			GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/np2p-alpine-linux-amd64 .
					GOOS=linux GOARCH=386 CGO_ENABLED=1 go build -o bin/np2p-gnu-linux-386 .
	CC=aarch64-linux-gnu-gcc 	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -o bin/np2p-gnu-linux-arm64 .
	CC=aarch64-linux-musl-gcc 	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -o bin/np2p-alpine-linux-arm64 .

archlinux-deps-prod: archlinux-deps-init archlinux-deps

archlinux-deps-init:
	sudo pacman -Syu --needed --noconfirm git base-devel
	git clone https://aur.archlinux.org/yay-bin.git
	mkdir -p ~/Downloads
	cd ~/Downloads
	cd yay-bin
	makepkg -si
	make arch-linux-deps
	cd -
archlinux-deps:
	yay --noconfirm --needed -S aarch64-linux-gnu-gcc musl aarch64-linux-musl go
