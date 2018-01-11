# build phase - see also:
# https://medium.com/travis-on-docker/multi-stage-docker-builds-for-creating-tiny-go-images-e0e1867efe5a
# https://medium.com/travis-on-docker/triple-stage-docker-builds-with-go-and-angular-1b7d2006cb88

# docker build -t intd-server .
# docker run -it -p 8080:8080 intd-server

FROM golang:alpine AS build-env

RUN apk add --update alpine-sdk

ADD . /go-artisanal-integers

RUN cd /go-artisanal-integers; make bin

FROM alpine

RUN apk add --update mysql 

COPY --from=build-env /go-artisanal-integers/bin/intd-server /intd-server

EXPOSE 8080

CMD /intd-server -host ${HOST} -port ${PORT} -db ${DB} -dsn ${DSN} -last ${LAST} -offset ${OFFSET} -increment ${INCREMENT}
