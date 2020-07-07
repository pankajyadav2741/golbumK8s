FROM golang:1.14.1-alpine3.11 AS base

LABEL maintainer="Pankaj Yadav <pankajyadav2741@gmail.com>"

WORKDIR /app

RUN apk update -qq && apk add git

RUN go get github.com/gocql/gocql && \
    go get github.com/gorilla/mux && \
	go get github.com/kelseyhightower/envconfig

COPY . .

RUN go build -o main .

FROM alpine

WORKDIR /album

COPY --from=base /app/ .

EXPOSE 5000

CMD [ "./main" ]
