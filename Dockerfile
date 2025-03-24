FROM alpine:3.21 AS health-downloader
ENV GRPC_HEALTH_PROBE_VERSION=v0.4.37 \
    GRPC_HEALTH_PROBE_URL=https://github.com/grpc-ecosystem/grpc-health-probe/releases/download
RUN apk -U add curl \
 && curl -fLso /bin/grpc_health_probe \
    ${GRPC_HEALTH_PROBE_URL}/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 \
 && chmod +x /bin/grpc_health_probe

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=health-downloader /bin/grpc_health_probe /bin/grpc_health_probe
COPY bin/server /masterdata-api
CMD ["/masterdata-api"]
