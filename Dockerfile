FROM golang:1.14.1-alpine3.11 AS base

LABEL maintainer="Pankaj Yadav <pankajyadav2741@gmail.com>"

WORKDIR /go/src/github.com/pankajyadav2741/golbumK8s/

COPY . .

RUN apk update -qq && apk add git

RUN go get github.com/gocql/gocql && \
    go get github.com/gorilla/mux && \
	go get github.com/kelseyhightower/envconfig 

RUN go build -o main .

FROM alpine

WORKDIR /go/src/github.com/pankajyadav2741/golbumK8s/

COPY --from=base /go/src/github.com/pankajyadav2741/golbumK8s/ .

EXPOSE 5000

CMD [ "./main" ]
