FROM golang:latest

WORKDIR /goprojects/MadTest

ADD ./ ./

RUN go build -o main .

EXPOSE 90

CMD [ "./main" ]
