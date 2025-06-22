FROM alpine:3.22
WORKDIR /app
COPY backupman .
EXPOSE 8080
ENV GIN_MODE=release
ENTRYPOINT ["/app/backupman"]
