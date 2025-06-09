FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .env . 

RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/wait-for-postgres.sh .
COPY --from=builder /app/db/migrations ./db/migrations

EXPOSE 8080

CMD ["./main"]