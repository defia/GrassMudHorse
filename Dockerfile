# dedicatedproxy on docker
FROM golang:1.5.1
MAINTAINER defia <defiaq@gmail.com>

COPY * $GOPATH/src/app/
WORKDIR $GOPATH/src/app/
EXPOSE 8888
RUN go get && go install
ENTRYPOINT wget -o=config.json $URL && $GOPATH/bin/app
