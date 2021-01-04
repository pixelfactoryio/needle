FROM golang:1.15-alpine AS builder

RUN apk --no-cache update && \
    apk --no-cache upgrade && \
    apk --no-cache add git build-base

WORKDIR /build

COPY . .

RUN go mod download

RUN make bin/needle

FROM alpine:3.12

COPY --from=builder /build/bin/needle /usr/bin/needle

ENTRYPOINT ["/usr/bin/needle"]
