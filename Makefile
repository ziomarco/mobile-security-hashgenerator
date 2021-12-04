.DEFAULT_GOAL := build

build:
	make clean
	@echo "Building..."
	go build
	GOOS=windows GOARCH=amd64 go build -o dist/msh-amd64.exe main.go
	GOOS=windows GOARCH=386 go build -o dist/msh-386.exe main.go
	GOOS=darwin GOARCH=amd64 go build -o dist/msh-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o dist/msh-darwin-arm64 main.go
	GOOS=linux GOARCH=386 go build -o dist/msh-linux-386 main.go
	GOOS=linux GOARCH=amd64 go build -o dist/msh-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o dist/msh-linux-arm64 main.go
	rm -rf mobile-security-hashgenerator
clean:
	@echo "Cleaning up"
	rm -rfv dist