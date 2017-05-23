FROM golang:1.8-alpine

ADD ./ /usr/src/envreplace
WORKDIR /usr/src/envreplace
RUN apk add --update make
RUN make install
