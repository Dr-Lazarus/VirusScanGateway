FROM golang:latest

WORKDIR /app

# Copy the application files
COPY . .


RUN go mod download


EXPOSE 8080 80 443


COPY entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["entrypoint.sh"]
CMD ["go", "run", "./cmd/server/main.go"]
