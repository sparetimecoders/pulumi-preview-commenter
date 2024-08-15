FROM alpine:20240807 AS build

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY pulumi-preview-commenter /
ENTRYPOINT ["/pulumi-preview-commenter"]
