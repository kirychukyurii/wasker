Kind = "ingress-gateway"
Name = "ingress-gateway"

Listeners = [
  {
    Port = 8443
    Protocol = "http2"
    Services = [
      {
        Name = "wasker-directory"
        Hosts = [
          "*"
        ]
      }
    ]
  }
]
