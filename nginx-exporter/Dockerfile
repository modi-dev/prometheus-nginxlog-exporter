FROM docker.repo-ci.sfera.inno.local/cont-docker-pub/registry.astralinux.ru/library/astra/ubi18-golang121:1.8.5 AS builder

ARG EXPORTER_VERSION=themkarimi:chore/update-dependencies

RUN mkdir -p /out
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOBIN=/out go install -trimpath -ldflags="-s -w" github.com/martin-helmich/prometheus-nginxlog-exporter@${EXPORTER_VERSION}

FROM docker.repo-ci.sfera.inno.local/cont-docker-pub/registry.astralinux.ru/library/astra/ubi18:1.8.5

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata logrotate \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/prometheus-nginxlog-exporter /usr/local/bin/prometheus-nginxlog-exporter

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
COPY logrotate/access-logrotate.conf /etc/logrotate.d/access-log

RUN chmod +x /usr/local/bin/docker-entrypoint.sh /usr/local/bin/prometheus-nginxlog-exporter && chmod 644 /etc/logrotate.d/access-log

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["-config-file", "/etc/prometheus-nginxlog-exporter/config.hcl"]
