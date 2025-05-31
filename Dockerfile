FROM apline:3.22
COPY backupman /
EXPOSE 8080
ENTRYPOINT ["/backupman"]
