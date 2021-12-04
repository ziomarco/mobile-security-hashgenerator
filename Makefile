.DEFAULT_GOAL := build

build:
	make clean
	@echo "Building..."
	go build
	GOOS=windows GOARCH=amd64 go build -o dist/msh-amd64.exe main.go
	GOOS=windows GOARCH=386 go build -o dist/msh-386.exe main.go
	GOOS=darwin GOARCH=amd64 go build -o dist/msh-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o dist/msh-arm64 main.go
clean:
	@echo "Cleaning up"
	rm -rfv dist