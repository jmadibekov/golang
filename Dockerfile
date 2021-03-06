# syntax=docker/dockerfile:1

# Awesome tutorial on containerizing Go application:
# https://docs.docker.com/language/golang/
FROM golang:1.17.1-alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY ./ ./

# Build
RUN go build -o /lostify ./project/cmd/main.go

# This is for documentation purposes only.
# To actually open the port, runtime parameters
# must be supplied to the docker command.
EXPOSE 8080

# Environment variable that our dockerised
# application can make use of. The value of environment
# variables can also be set via parameters supplied
# to the docker command on the command line.
ENV SERVER_TYPE=http

# Run
CMD [ "/lostify" ]
