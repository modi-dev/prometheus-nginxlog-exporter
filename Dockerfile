FROM registry.astralinux.ru/library/astra/ubi18-golang121:1.8.5 AS builder

ARG EXPORTER_VERSION=themkarimi:chore/update-dependencies

COPY . /out/
RUN cd /out && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOBIN=/out go install

FROM registry.astralinux.ru/library/astra/ubi18:1.8.5

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata logrotate \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/prometheus-nginxlog-exporter /usr/local/bin/prometheus-nginxlog-exporter

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
COPY logrotate/access-logrotate.conf /etc/logrotate.d/access-log

RUN chmod +x /usr/local/bin/docker-entrypoint.sh /usr/local/bin/prometheus-nginxlog-exporter && chmod 644 /etc/logrotate.d/access-log

# Create a non-root user, set permissions, and copy application files in a single RUN instruction
RUN groupadd -r appuser && \
    useradd -r -g appuser appuser && \
    mkdir -p /opt/app && \
    chown -R appuser:appuser /opt/app
# Switch to the non-root user
USER appuser
 
# Add a health check
HEALTHCHECK --interval=30s --timeout=10s --retries=5 \
  CMD nc -z localhost 9113 || exit 1

# Expose ports
EXPOSE 9113/tcp

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["-config-file", "/etc/prometheus-nginxlog-exporter/config.hcl"]
