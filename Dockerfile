# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
# FROM golang:1.13.1-alpine

# Set the Current Working Directory inside the container
# WORKDIR /app

# Copy go mod and sum files
# COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
# COPY . .

# Build the Go app
# RUN make build

# Command to run the executable
# CMD ["./go-play-publisher"]


FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go build -o main .

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]
