# Stage 1: Build
FROM golang:1-bookworm

LABEL maintainer="srad <ssedigi@posteo.de>"
LABEL source="https://github.com/srad/wordpress-backup-enhanced"

RUN apt update && apt upgrade -y
RUN apt install default-mysql-client bzip2 -y
RUN which mysqldump

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

ENV MYSQL_HOST=""
ENV MYSQL_USER=""
ENV MYSQL_PASSWORD=""
ENV MYSQL_DATABASE=""
ENV BASIC_AUTH_PASSWORD=""

ENV MYSQL_HOST_FILE=""
ENV MYSQL_USER_FILE=""
ENV MYSQL_PASSWORD_FILE=""
ENV MYSQL_DATABASE_FILE=""
ENV BASIC_AUTH_PASSWORD_FILE=""

ENV PORT="8080"

ENV BACKUP_ON_START="true"

ENV BACKUP_CRON="@every 1h"

RUN mkdir /backups
VOLUME /backups

EXPOSE $PORT

ENTRYPOINT ["/docker-entrypoint.sh"]
