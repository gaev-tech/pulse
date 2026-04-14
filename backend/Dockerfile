FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .
RUN swag init -g cmd/api/main.go -o api/docs
RUN CGO_ENABLED=0 GOOS=linux go build -o /pulse ./cmd/api

FROM golang:1.25-alpine

WORKDIR /app
COPY --from=builder /pulse /app/pulse

EXPOSE 8080

CMD ["/app/pulse"]
