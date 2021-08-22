FROM alpine:latest

COPY needle_*.apk /tmp/

RUN apk add --allow-untrusted /tmp/needle_*.apk \
    && rm -fr /tmp/needle_*.apk

ENTRYPOINT ["/usr/local/bin/needle"]
