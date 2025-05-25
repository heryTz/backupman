build:
	rm -rf ./bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/backupman-linux-amd64 .
	GCGO_ENABLED=0 OOS=darwin GOARCH=amd64 go build -o ./bin/backupman-darwin-amd64 .
