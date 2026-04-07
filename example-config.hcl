{{- $service := .Release.Name }}

listen {
  port = 9113
}

namespace "nginx" {
  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\""
  source {
    files = [
      "/var/log/nginx/access.log"
    ]
  }
  relabel "request_uri" {
    from = "request"
    split = 2
    separator = " " // (1)

    match "^/+([^/]*).*" {
      replacement = "/$1"
    }
  }
  labels {
    app = "{{ $service }}"
  }
}