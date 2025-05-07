FROM ubuntu
WORKDIR /app
RUN apt update \
  && apt install -y systemd
COPY backupman-linux-amd64 /app/backupman
COPY backupman.service /etc/systemd/system/backupman.service
RUN systemctl enable backupman.service \
  && systemctl start backupman.service
EXPOSE 8080
