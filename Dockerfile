FROM golang:1.12-alpine as builder

WORKDIR /app
COPY . .

RUN apk add --no-cache \
    gcc \
    git \
    linux-headers \
    make \
    musl-dev

RUN make build

FROM alpine:3.9

COPY --from=builder /app/bin/ /usr/local/bin/

RUN addgroup -g 1000 sauce && \
    adduser -h /sauce -D -u 1000 -G sauce sauce && \
    chown sauce:sauce /sauce

USER sauce

ENTRYPOINT ["chef-scheduler"]
