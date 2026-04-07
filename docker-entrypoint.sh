#!/bin/sh
set -eu

LOGROTATE_STATE_FILE="${LOGROTATE_STATE_FILE:-/tmp/logrotate.status}"
LOGROTATE_INTERVAL_SECONDS="${LOGROTATE_INTERVAL_SECONDS:-60}"

touch "${LOGROTATE_STATE_FILE:-/tmp/logrotate.status}"

prometheus-nginxlog-exporter "$@" &
EXPORTER_PID=$!

term_handler() {
  kill -TERM "$EXPORTER_PID" 2>/dev/null || true
  wait "$EXPORTER_PID" 2>/dev/null || true
  exit 0
}

trap term_handler INT TERM

while true; do
  logrotate -s "$LOGROTATE_STATE_FILE" /etc/logrotate.d/access-log || true
  if ! kill -0 "$EXPORTER_PID" 2>/dev/null; then
    wait "$EXPORTER_PID"
    exit $?
  fi
  sleep "$LOGROTATE_INTERVAL_SECONDS"
done
