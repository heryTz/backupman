docker:
	docker build . --tag backupman
	docker run --rm \
		--publish 8080:8080 \
		--name backupman \
		backupman
