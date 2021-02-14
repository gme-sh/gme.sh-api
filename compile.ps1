New-Item -Path . -Name "bin" -ItemType "directory" 
Write-Host "Compiling for every OS and Platform"
Write-Host "🐧 Compile for Linux"
Set-Variable GOOS=linux 
Set-Variable GOARCH=amd64 
go build -o ./bin/gme-linux-amd64 ./cmd/gme-shortener/main.go
Set-Variable GOOS=linux 
Set-Variable GOARCH=386 
go build -o ./bin/gme-linux-386 ./cmd/gme-shortener/main.go  
Set-Variable GOOS=linux
Set-Variable GOARCH=arm 
go build -o ./bin/gme-linux-arm ./cmd/gme-shortener/main.go 
Set-Variable GOOS=linux
Set-Variable GOARCH=arm64 
go build -o ./bin/gme-linux-arm64 ./cmd/gme-shortener/main.go
Write-Host "🍏 Compile for Apple"
Set-Variable GOOS=darwin 
Set-Variable GOARCH=amd64 
go build -o ./bin/gme-darwin-amd64 ./cmd/gme-shortener/main.go
Write-Host "🪟 Compile for Windows"
Set-Variable GOOS=windows 
Set-Variable GOARCH=amd64 
go build -o ./bin/gme-windows-amd64 ./cmd/gme-shortener/main.go
Set-Variable GOOS=windows 
Set-Variable GOARCH=386 
go build -o ./bin/gme-windows-386 ./cmd/gme-shortener/main.go
Write-Host "🐡 Compile for FreeBSD"
Set-Variable GOOS=freebsd 
Set-Variable GOARCH=amd64 
go build -o ./bin/gme-freebsd-amd64 ./cmd/gme-shortener/main.go
Set-Variable GOOS=freebsd 
Set-Variable GOARCH=386 
go build -o ./bin/gme-freebsd-386 ./cmd/gme-shortener/main.go  
Set-Variable GOOS=freebsd 
Set-Variable GOARCH=arm 
go build -o ./bin/gme-freebsd-arm ./cmd/gme-shortener/main.go