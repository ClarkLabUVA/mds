FROM golang:latest
RUN go get go.mongodb.org/mongo-driver/mongo \
  go.mongodb.org/mongo-driver/bson \
  github.com/satori/go.uuid \
  github.com/gorilla/mux \
  github.com/urfave/negroni \
  github.com/dgrijalva/jwt-go

WORKDIR /mds
COPY . /go/src/mds

RUN go build -o main mds

EXPOSE 80

ENTRYPOINT ["./main"]
