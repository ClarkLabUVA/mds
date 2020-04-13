FROM golang:latest
#FROM ruby:2.6
#RUN gem install bolognese
#COPY bin/go1.14.linux-amd64.tar.gz .
#RUN tar -xf go1.14.linux-amd64.tar.gz && rm go1.14.linux-amd64.tar.gz

ENV PATH=$PATH:/go/bin
ENV GOPATH=/usr/go

RUN go get go.mongodb.org/mongo-driver/mongo \
 go.mongodb.org/mongo-driver/bson \
 github.com/satori/go.uuid \
 github.com/gorilla/mux \
 github.com/urfave/negroni 

WORKDIR /mds
COPY src/ .

RUN go build .

EXPOSE 80

ENTRYPOINT ["./main"]
