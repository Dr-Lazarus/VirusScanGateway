FROM golang:latest

WORKDIR /app

# Copy the application files
COPY . .

# Download Go modules
RUN go mod download

# Expose necessary ports
EXPOSE 8080 80 443

# Use an entrypoint script to decide at runtime what to do
COPY entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["entrypoint.sh"]
CMD ["go", "run", "./cmd/server/main.go"]
