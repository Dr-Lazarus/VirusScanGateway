# Use the official Golang image.
FROM golang:latest

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /app

# Copy the local package files to the container's workspace.
COPY . .

# Download all dependencies.
RUN go mod download

# Expose port 8080 to the outside world.
EXPOSE 8080

# Run the Go application directly with `go run`.
CMD ["go", "run", "./cmd/server/main.go"]
