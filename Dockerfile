FROM scratch
COPY backupman /
ENTRYPOINT ["/backupman"]
