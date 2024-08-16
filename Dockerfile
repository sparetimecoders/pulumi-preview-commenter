FROM alpine:20240807@sha256:b93f4f6834d5c6849d859a4c07cc88f5a7d8ce5fb8d2e72940d8edd8be343c04 AS build

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY pulumi-preview-commenter /
ENTRYPOINT ["/pulumi-preview-commenter"]
