# syntax=docker/dockerfile:1


FROM golang:1.17.5-alpine

WORKDIR /app

COPY . ./

RUN go build -o /jokesapp
EXPOSE 9090
CMD [ "/jokesapp" ]
