# build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bitedash ./cmd/api

# runtime stage
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bitedash .

EXPOSE 8080
EXPOSE 9090

CMD ["./bitedash"]