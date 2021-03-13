FROM golang:latest AS build
WORKDIR /go/src/app
ENV CGO_ENABLED=0

# Install dependencies
# Thanks to @montanaflynn
# https://github.com/montanaflynn/golang-docker-cache
COPY go.mod go.sum ./
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get

# Copy source-code and build
COPY . .
RUN go build ./cmd/gme-sh


LABEL maintainer="darmiel <hi@d2a.io>"


FROM scratch
COPY --from=build /go/src/app/gme-sh /
CMD ["/gme-sh"]