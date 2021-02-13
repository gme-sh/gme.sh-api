FROM golang:latest

LABEL maintainer="darmiel <hi@d2a.io>"

WORKDIR /usr/app/src
COPY . .

RUN go get ./cmd/gme-shortener
RUN go build ./cmd/gme-shortener

CMD ["./gme-shortener"]