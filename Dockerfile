FROM alpine:3.13.3

RUN apk add --no-cache ca-certificates

ADD ./cert-operator /cert-operator

ENTRYPOINT ["/cert-operator"]
