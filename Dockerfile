FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main ./main.go

COPY xplore-48-447519269b91.json /app/credentials/xplore-48-447519269b91.json


EXPOSE 8080

CMD ["./main"]
