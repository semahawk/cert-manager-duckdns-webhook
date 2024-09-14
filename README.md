<p align="center">
  <img src="https://raw.githubusercontent.com/cert-manager/cert-manager/d53c0b9270f8cd90d908460d69502694e1838f5f/logo/logo-small.png" height="256" width="256" alt="cert-manager project logo" />
  &nbsp;
  <img src="https://raw.githubusercontent.com/linuxserver/docker-templates/master/linuxserver.io/img/duckdns.png" height="256" width="256" alt="DuckDNS project logo" />
</p>

# DuckDNS Webhook for cert-manager

This WebHook solves the DNS01 challenge to prove ownership of DuckDNS domains.

## Helm Chart

TODO

### Running the test suite

1. Create a domain in DuckDNS.
2. Run the following command.
    ```shell
    export DUCKDNS_DOMAIN=<your-domain> // Example: my-domain for my-domain.duckdns.org
    export DUCKDNS_TOKEN=<your-token>
    ```
3. Run `make test`

### Pushing the image

```shell
  export IMAGE_TAG=1.0.0
  docker buildx build --push --platform linux/arm64,linux/amd64 -t csp33/cert-manager-duckdns-webhook:$IMAGE_TAG -t csp33/cert-manager-duckdns-webhook:latest .;
```

## Acknowledgments

- Inspired by [the wonderful ebrianne's job](https://github.com/ebrianne/cert-manager-webhook-duckdns).