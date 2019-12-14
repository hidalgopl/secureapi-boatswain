FROM golang:1.12-alpine3.9 as builder

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
COPY --from=builder /app/config.yaml /etc/sailor/config.yaml

RUN addgroup -g 1000 boatswain && \
    adduser -h /boatswain -D -u 1000 -G boatswain boatswain

USER boatswain

ENTRYPOINT ["/usr/local/bin/boatswain",  "run"]
