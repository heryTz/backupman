FROM ubuntu
WORKDIR /app
RUN apt update \
  && apt install -y supervisor
COPY backupman-linux-amd64 /app/backupman
COPY supervisord.conf /etc/supervisord.conf
EXPOSE 8080
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisord.conf"]
