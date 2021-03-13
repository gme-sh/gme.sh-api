FROM golang:latest
WORKDIR /go/src/app
LABEL maintainer="darmiel <hi@d2a.io>"

# Install dependencies
# Thanks to @montanaflynn
# https://github.com/montanaflynn/golang-docker-cache
COPY go.mod go.sum ./
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get

# Copy source-code and build
COPY . .
RUN go build ./cmd/gme-sh

CMD ["./gme-sh"]