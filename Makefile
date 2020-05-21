.PHONY: build-mac
build-mac:
	GOOS=darwin GOARCH=amd64 go build -o build/darwin/jp-name main.go

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o build/linux/jp-name main.go

.PHONY: clean
clean:
	rm -rf build/
