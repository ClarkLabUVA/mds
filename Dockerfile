FROM golang:latest as builder
WORKDIR /mds
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY go.mod .
RUN CGO_ENABLED=0 GOOS=linux go build -o mds cmd/mds/main.go

# required to run bolognese
#FROM ruby:2.6
#RUN gem install bolognese
#COPY bin/go1.14.linux-amd64.tar.gz .
#RUN tar -xf go1.14.linux-amd64.tar.gz && rm go1.14.linux-amd64.tar.gz


# runner image 
FROM alpine:latest
WORKDIR /mds
COPY --from=builder /mds/mds .
ENTRYPOINT ["./mds"]
