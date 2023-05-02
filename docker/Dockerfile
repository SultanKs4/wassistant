#Build
FROM golang:1.20-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -o wassistant cmd/wassistant/main.go

#Deploy
FROM alpine:latest

COPY --from=builder /app/wassistant /app/

COPY --from=builder /app/.env.docker /app/.env

CMD ["/app/wassistant"]