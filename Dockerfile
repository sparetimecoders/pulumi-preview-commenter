FROM alpine:20240807@sha256:931d2b47d03ca687b4306020ddac298ee75f1539ab2767049450b99c872e81d0 AS build

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY pulumi-preview-commenter /
ENTRYPOINT ["/pulumi-preview-commenter"]
