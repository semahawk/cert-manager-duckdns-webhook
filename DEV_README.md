# Running the test suite

1. Create a domain in DuckDNS.
2. Run the following command.
    ```shell
    export DUCKDNS_DOMAIN=<your-domain> // Example: my-domain for my-domain.duckdns.org
    export DUCKDNS_TOKEN=<your-token>
    ```
3. Run `make test`

# Pushing the image

```shell
export IMAGE_TAG=1.0.1
docker buildx build --push --platform linux/arm64,linux/amd64 -t csp33/cert-manager-duckdns-webhook:$IMAGE_TAG -t csp33/cert-manager-duckdns-webhook:latest .;
```
