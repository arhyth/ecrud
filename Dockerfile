# Use the official Golang base image
FROM golang:1.21-bookworm

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project to the working directory
COPY . .

# Build the Go server binary
RUN go build -o server ./cmd

# Expose the port the server will run on
EXPOSE 3000

# Set the entry point for the container
CMD ["./server"]
