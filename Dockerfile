FROM golang:latest

WORKDIR /app

COPY . .

COPY /home/ec2-user/VirusScanGateway/certs/virusscanapi.lat-ssl-bundle /etc/ssl/certs

RUN go mod download

EXPOSE 8080

CMD ["go", "run", "./cmd/server/main.go"]
