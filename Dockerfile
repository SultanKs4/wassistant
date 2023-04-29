#Build
FROM golang:1.20-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build cmd/wassistant/main.go

#Deploy
FROM alpine:latest

COPY --from=builder /app/cmd/wassistant /app/cmd/wassistant

CMD ["/app/cmd/wassistant"]