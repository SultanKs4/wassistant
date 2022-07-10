#Build
FROM golang:1.18.3-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build

#Deploy
FROM alpine:latest

COPY --from=builder /app/wassistant /app/wassistant

CMD ["/app/wassistant"]