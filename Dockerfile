FROM alpine:3.14.1

RUN apk add --no-cache ca-certificates

ADD ./cert-operator /cert-operator

ENTRYPOINT ["/cert-operator"]
