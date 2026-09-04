.PHONY: build install test clean release

BIN = mark
MAIN = ./cmd/mark

build:
	go build -ldflags="-s -w" -o $(BIN).exe $(MAIN)

install:
	go install -ldflags="-s -w" $(MAIN)

test:
	go test ./...

clean:
	del /f $(BIN).exe 2>nul || true

release:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/mark-windows-amd64.exe $(MAIN)
	GOOS=windows GOARCH=386   go build -ldflags="-s -w" -o dist/mark-windows-386.exe   $(MAIN)
	GOOS=linux   GOARCH=amd64 go build -ldflags="-s -w" -o dist/mark-linux-amd64       $(MAIN)
	GOOS=darwin  GOARCH=arm64 go build -ldflags="-s -w" -o dist/mark-darwin-arm64      $(MAIN)
