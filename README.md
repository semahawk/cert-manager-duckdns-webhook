<p align="center">
  <img src="https://raw.githubusercontent.com/cert-manager/cert-manager/d53c0b9270f8cd90d908460d69502694e1838f5f/logo/logo-small.png" height="256" width="256" alt="cert-manager project logo" />
  &nbsp;
  <img src="https://raw.githubusercontent.com/linuxserver/docker-templates/master/linuxserver.io/img/duckdns.png" height="256" width="256" alt="DuckDNS project logo" />
</p>

# DuckDNS Webhook for cert-manager

This WebHook solves the DNS01 challenge to prove ownership of DuckDNS domains.

## How to use it
1. Add the repository
   ```shell
   helm repo add csp33 https://csp33.github.io/cert-manager-duckdns-webhook
   ```
2. Create a values.yaml file with the following content:
   ```yaml
    token:
      value: <your DuckDNS token>
    clusterIssuer:
      email: <your email>
      production:
        create: true
      staging:
        create: true
   ```
3. Install the chart:
   ```shell
   helm install cert-manager-duckdns-webhook csp33/cert-manager-duckdns-webhook -f values.yaml
   ```
4. Add the following annotation to the ingress you want to generate a certificate for:
   ```yaml
   cert-manager.io/cluster-issuer: duckdns-letsencrypt-prod
   ```
5. Wait for it to finish


## Acknowledgments

- Inspired by [the wonderful ebrianne's job](https://github.com/ebrianne/cert-manager-webhook-duckdns).
