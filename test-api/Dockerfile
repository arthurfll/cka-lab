FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -o api main.go

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/api .

RUN mkdir -p /app/data

EXPOSE 7000

CMD ["./api"]