compile:
	mkdir -p bin/
	@echo "Compiling for every OS and Platform"
	@echo "🐧 Compile for Linux"
	GOOS=linux GOARCH=amd64 go build -o ./bin/gme-linux-amd64 ./cmd/gme-sh/main.go
	GOOS=linux GOARCH=386 go build -o ./bin/gme-linux-386 ./cmd/gme-sh/main.go
	GOOS=linux GOARCH=arm go build -o ./bin/gme-linux-arm ./cmd/gme-sh/main.go
	GOOS=linux GOARCH=arm64 go build -o ./bin/gme-linux-arm64 ./cmd/gme-sh/main.go
	@echo "🍏 Compile for Apple"
	GOOS=darwin GOARCH=amd64 go build -o ./bin/gme-darwin-amd64 ./cmd/gme-sh/main.go
	@echo "🪟 Compile for Windows"
	GOOS=windows GOARCH=amd64 go build -o ./bin/gme-windows-amd64 ./cmd/gme-sh/main.go
	GOOS=windows GOARCH=386 go build -o ./bin/gme-windows-386 ./cmd/gme-sh/main.go
	@echo "🐡 Compile for FreeBSD"
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/gme-freebsd-amd64 ./cmd/gme-sh/main.go
	GOOS=freebsd GOARCH=386 go build -o ./bin/gme-freebsd-386 ./cmd/gme-sh/main.go
	GOOS=freebsd GOARCH=arm go build -o ./bin/gme-freebsd-arm ./cmd/gme-sh/main.go