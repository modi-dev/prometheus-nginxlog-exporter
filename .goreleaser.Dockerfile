FROM scratch

COPY prometheus-nginxlog-exporter /prometheus-nginxlog-exporter

EXPOSE 9113
ENTRYPOINT ["/prometheus-nginxlog-exporter"]
