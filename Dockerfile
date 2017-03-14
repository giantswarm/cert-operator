FROM alpine:3.5

RUN apk add --update ca-certificates \
    && rm -rf /var/cache/apk/*

ADD ./cert-operator /cert-operator

ENTRYPOINT ["/cert-operator"]
