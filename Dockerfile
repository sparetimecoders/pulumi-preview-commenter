FROM alpine:20240807@sha256:8431297eedca8df8f1e6144803c6d7e057ecff2408aa6861213cb9e507acadf8 AS build

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY pulumi-preview-commenter /
ENTRYPOINT ["/pulumi-preview-commenter"]
