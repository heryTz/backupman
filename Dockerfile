FROM alpine:3.22
WORKDIR /app
COPY backupman .
EXPOSE 8080
ENTRYPOINT ["/app/backupman"]
