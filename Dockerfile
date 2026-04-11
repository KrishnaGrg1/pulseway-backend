FROM golang:1.25-alpine AS BUILDER

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

RUN CGO_ENABLED=0 go build -o pulseway ./cmd
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --FROM=BUILDER /app/pulseway

EXPOSE 8080

CMD ["./pulseway"]