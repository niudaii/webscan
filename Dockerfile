FROM golang:1.18.4-alpine as builder
RUN apk add --no-cache make git
WORKDIR /webscan-src
COPY . /webscan-src
RUN go mod download && \
    make docker && \
    mv ./bin/webscan-docker /webscan

FROM alpine:latest
COPY --from=builder /webscan /

ENTRYPOINT ["/webscan"]