FROM alpine:20240923@sha256:f4b9f111e2c5290552a920590dd48dc58f5ea1cacda6e25b0a2718974d090cf0 AS build

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY pulumi-preview-commenter /
ENTRYPOINT ["/pulumi-preview-commenter"]
