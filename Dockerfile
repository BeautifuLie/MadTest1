# syntax=docker/dockerfile:1


FROM golang:1.17.5-alpine

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

COPY *.go ./

RUN go mod download

COPY . ./

RUN go build -o /jokesapp
EXPOSE 9090
CMD [ "/jokesapp" ]
