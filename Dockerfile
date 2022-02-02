# syntax=docker/dockerfile:1


FROM golang:1.17.5-alpine AS builder

WORKDIR /app

COPY . ./

RUN go build -o jokesapp .


FROM alpine:latest
WORKDIR /myapp
COPY --from=builder /app /myapp/

EXPOSE 9090
CMD ["./jokesapp"]
