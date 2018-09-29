FROM golang:latest

ADD . /go/src/github.com/KNaiskes/Home-Panel

RUN go get github.com/mattn/go-sqlite3 \
    && go get github.com/eclipse/paho.mqtt.golang \
    && go get github.com/gorilla/sessions

RUN go install github.com/KNaiskes/Home-Panel
ENTRYPOINT /go/bin/Home-Panel

EXPOSE 8080
