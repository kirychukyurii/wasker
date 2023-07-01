service {
  kind = "ingress-gateway"

  id = "ingress-gateway"
  name = "ingress-gateway"

  address = "127.0.0.1"
  port = 8443

  checks = [
    {
      id       = "envoy_ready_listener"
      name     = "HTTP on :8444/ready"
      http     = "http://127.0.0.1:8444/ready"
      interval = "15s"
    },
    {
      id       = "envoy_public_listener"
      name     = "TCP on :8000"
      tcp      = "127.0.0.1:8000"
      interval = "15s"
    },
  ]

  proxy = {
    config = {
      protocol = "http2"
    }
  }
}