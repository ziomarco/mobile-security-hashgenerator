.DEFAULT_GOAL := build

build:
	@echo "Building..."
	go build
	mkdir -p dist
	mv mobile-security-hashgenerator dist/msh

clean:
	@echo "Cleaning up"
	rm -rfv dist