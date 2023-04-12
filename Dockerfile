FROM alpine:3.17 as health-downloader
ENV GRPC_HEALTH_PROBE_VERSION=v0.4.17 \
    GRPC_HEALTH_PROBE_URL=https://github.com/grpc-ecosystem/grpc-health-probe/releases/download
RUN apk -U add curl \
 && curl -fLso /bin/grpc_health_probe \
    ${GRPC_HEALTH_PROBE_URL}/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 \
 && chmod +x /bin/grpc_health_probe

FROM ghcr.io/metal-stack/builder:latest as builder

FROM alpine:3.17
RUN apk -U add ca-certificates
COPY --from=builder /work/bin/server /masterdata-api
COPY --from=health-downloader /bin/grpc_health_probe /bin/grpc_health_probe
CMD ["/masterdata-api"]
