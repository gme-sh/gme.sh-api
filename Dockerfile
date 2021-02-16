FROM golang:latest

LABEL maintainer="darmiel <hi@d2a.io>"

WORKDIR /usr/app/src
COPY . .

RUN go get ./cmd/gme-sh
RUN go build ./cmd/gme-sh

CMD ["./gme-sh"]