FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -ldflags="-w -s" -a -installsuffix cgo -o main ./main.go

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /build/main /app/main

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/main"]
