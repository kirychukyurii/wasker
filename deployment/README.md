# Deployments
## Consul Service Mesh + Envoy Proxy
### Config
Genarate Envoy Ingress Proxy config (with proto file descriptors)

    $ ./transflect --plaintext --format envoy --http-port 8443 localhost:8086 envoy.yaml

Apply Consul Service Defaults definitions

    $ consul config write deployment/consul/defaults/service-defaults.hcl

Apply Consul Proxy Defaults definitions

    $ consul config write deployment/consul/defaults/proxy-defaults.hcl

Apply Consul Ingress Listeners

    $ consul config write deployment/consul/ingress-gateway/listeners.hcl

Register Consul Service for Envoy Ingress Gateway

    $ consul services register deployment/consul/ingress-gateway/service.hcl

### Run
Start Envoy Ingress Proxy

    $ envoy -c deployment/envoy/envoy.yml
