FROM alpine:3.12.2

RUN apk add --no-cache ca-certificates

ADD ./cert-operator /cert-operator

ENTRYPOINT ["/cert-operator"]
