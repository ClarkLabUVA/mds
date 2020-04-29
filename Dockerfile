FROM golang:latest
#FROM ruby:2.6
#RUN gem install bolognese
#COPY bin/go1.14.linux-amd64.tar.gz .
#RUN tar -xf go1.14.linux-amd64.tar.gz && rm go1.14.linux-amd64.tar.gz

WORKDIR /mds
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY go.mod .

RUN go build -o mds cmd/mds/main.go
EXPOSE 80

ENTRYPOINT ["./mds"]
