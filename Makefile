compile:
	@echo "Compiling for every OS and Platform"
	@echo "Compile for Linux"
	GOOS=linux GOARCH=amd64 go build -o ./bin/main-linux-amd64 ./cmd/gme-shortener/main.go
	GOOS=linux GOARCH=386 go build -o ./bin/main-linux-386 ./cmd/gme-shortener/main.go  
	GOOS=linux GOARCH=arm go build -o ./bin/main-linux-arm ./cmd/gme-shortener/main.go 
	GOOS=linux GOARCH=arm64 go build -o ./bin/main-linux-arm64 ./cmd/gme-shortener/main.go 
	@echo "Compile for Windows"
	GOOS=windows GOARCH=amd64 go build -o ./bin/main-windows-amd64 ./cmd/gme-shortener/main.go
	GOOS=windows GOARCH=386 go build -o ./bin/main-windows-386 ./cmd/gme-shortener/main.go  
	@echo "Compile for Apple"
	GOOS=darwin GOARCH=arm go build -o ./bin/main-darwin-arm ./cmd/gme-shortener/main.go 
	GOOS=darwin GOARCH=arm64 go build -o ./bin/main-darwin-arm64 ./cmd/gme-shortener/main.go
	GOOS=darwin GOARCH=amd64 go build -o ./bin/main-darwin-amd64 ./cmd/gme-shortener/main.go
	GOOS=darwin GOARCH=386 go build -o ./bin/main-darwin-386 ./cmd/gme-shortener/main.go   
	@echo "Compile for FreeBSD"
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/main-freebsd-amd64 ./cmd/gme-shortener/main.go
	GOOS=freebsd GOARCH=386 go build -o ./bin/main-freebsd-386 ./cmd/gme-shortener/main.go  
	GOOS=freebsd GOARCH=arm go build -o ./bin/main-freebsd-arm ./cmd/gme-shortener/main.go