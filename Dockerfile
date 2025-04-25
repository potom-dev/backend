# syntax=docker/dockerfile:1

FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o backend .

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/backend ./backend
EXPOSE 8080
ENTRYPOINT ["/app/backend"]