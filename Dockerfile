FROM golang:alpine AS builder

RUN apk --no-cache -U add git

RUN go get -u github.com/kardianos/govendor

WORKDIR /go/src/app
COPY ./Server .

RUN govendor sync

RUN govendor build -o /go/src/app/webserver


FROM alpine

COPY --from=builder /go/src/app/webserver /

CMD [ "/webserver" ]
