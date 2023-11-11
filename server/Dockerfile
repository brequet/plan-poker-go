# FROM node:20-slim AS base
# Use the official Golang image as the base image
FROM golang:alpine AS build

# Add the necessary build tools
RUN apk --no-cache add gcc g++ make git

# Set the Current Working Directory inside the container
WORKDIR /go/src/app

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

# Download all dependencies. 
# Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Build the Go app
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/web-app .

# Start a new stage from scratch
FROM alpine:3.13

# Add the certificates
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /usr/bin

# Copy the pre-built binary file from the previous stage
COPY --from=build /go/src/app/bin /go/bin

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
ENTRYPOINT /go/bin/web-app --port 8080