build:
	GOOS=linux GOARCH=amd64 go build -o backupman-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o backupman-darwin-amd64 main.go

docker:
	docker build . --tag backupman
		
	docker run -it --rm \
		--publish 8080:8080 \
		--name backupman \
		backupman \
		sh
