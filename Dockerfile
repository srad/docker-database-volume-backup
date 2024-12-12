FROM node:22 as nodebuild

WORKDIR /app
COPY ./frontend .
RUN npm install
RUN npm run build

FROM ubuntu:24.04

LABEL maintainer="srad <ssedighi@posteo.de>"
LABEL source="https://github.com/srad/docker-database-volume-backup"

RUN apt update && apt upgrade -y
RUN apt install mysql-client bzip2 wget -y

RUN wget -qO- https://go.dev/dl/go1.23.4.linux-amd64.tar.gz | tar xvz -C /usr/local
RUN export PATH=$PATH:/usr/local/go/bin

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN /usr/local/go/bin/go mod download

# Copy the rest of the application code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux /usr/local/go/bin/go build -o /main

COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

COPY --from=nodebuild /app/dist ./public

ENV MYSQL_HOST=""
ENV MYSQL_USER=""
ENV MYSQL_PASSWORD=""
ENV MYSQL_DATABASE=""
ENV BASIC_AUTH_USER=""
ENV BASIC_AUTH_PASSWORD=""

ENV MYSQL_HOST_FILE=""
ENV MYSQL_USER_FILE=""
ENV MYSQL_PASSWORD_FILE=""
ENV MYSQL_DATABASE_FILE=""
ENV BASIC_AUTH_USER_FILE=""
ENV BASIC_AUTH_PASSWORD_FILE=""

ENV PORT="8080"

ENV BACKUP_ON_START="true"
ENV BACKUP_CRON="@every 1h"

RUN mkdir /backups
RUN mkdir /backups/dumps
RUN mkdir /backups/files
VOLUME /backups

RUN mkdir /files
VOLUME /files

EXPOSE $PORT

ENTRYPOINT ["/docker-entrypoint.sh"]
