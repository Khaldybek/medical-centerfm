FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
RUN go build -o auth-service

EXPOSE 8080

CMD ["./auth-service"] 