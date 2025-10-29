FROM golang:1.21-alpine AS builder

WORKDIR /build

COPY . .

RUN go mod download
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o voice-server -v ./main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /build/voice-server .

EXPOSE 3002

CMD ["./voice-server"]
