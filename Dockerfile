FROM golang:latest

WORKDIR /app


COPY . .


RUN go mod download


EXPOSE 8080 80 443


CMD ["go", "run", "./cmd/server/main.go"]
