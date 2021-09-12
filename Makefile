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

#Go compiling using CGo(sqlite) AND for Alpine on top of that created most of this confusion

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

arch-container-init:
	echo 'ParallelDownloads=5' >> /etc/pacman.conf
	useradd -m -s /bin/sh builder
	pacman -Syu sudo --noconfirm --needed
	echo 'builder ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
	chmod ugo+rwx .
	sudo -u builder make archlinux-deps-prod build-prod
