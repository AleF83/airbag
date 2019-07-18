# airbag

Airbag is a tiny OAuth2 sidecar.

## Usage

To start to use Airbag add it to pod/deployment definition and configure it as described below.
Then set the service port to the airbag's one and that's it.

## Configuration

To configure Airbag you need provide to it details for JWT token validation. Airbag supports authentication by several issuers. So you need to provide list of `iss`/`aud` pairs. Airbag will verify that JWT is sign correctly and its claims (`iss`,`aud`) fit to one of these configurations.

> Note: In some cases OpenID Connect server hasn't discovery page or it's discovery page URL isn't in format: ***{issuer}/.well-known/openid-configuration***. In this case `jwks_url` property should be provided aside `iss` one.

Use the following env variables to configure Airbag:

- `AIRBAG_PORT`: Declares port for airbag (usually 443 or 80).
- `AIRBAG_CONFIG_ENV_PREFIX`: Airbag scan all env variables with this prefix for it's configuration. Not operational now. For future use.
- `AIRBAG_CONFIG_NAME`: Name of JSON file that contains JWT provider configurations.
- `AIRBAG_CONFIG_PATH`: Path to configuration file.

## Example

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-service
  namespace: my-namespace
spec:
  selector:
    matchLabels:
      app: my-service
  template:
    metadata:
      labels:
        app: my-service
    spec:
      containers:
      - name: my-service
        image: my-image
        livenessProbe:
          httpGet:
            path: /users
            port: 8080
          initialDelaySeconds: 3
          periodSeconds: 30
      - name: airbag
        image: alex4soluto/airbag
        env:
          - name: AIRBAG_PORT
            value: 443
          - name: AIRBAG_CONFIG_ENV_PREFIX
            value: AIRBAG
          - name: AIRBAG_CONFIG_NAME
            value: airbag-config
          - name: AIRBAG_CONFIG_PATH
            value: /app/data
        volumeMounts:
          - name: data
            mountPath: /app/data/airbag-config.json
            subPath: airbag-config.json
            readOnly: true
      volumes:
        - name: data
          configMap:
            name: my-service-config-map
---
kind: Service
apiVersion: v1
metadata:
  name: my-service
  namespace: my-namespace
spec:
  selector:
    app: my-service
  ports:
    - name: https
      port: 443
      targetPort: 443
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-service-config-map
  namespace: my-namespace
data:
  airbag-config.json: |
        {
          "backend": "http://localhost:8080",
          "JWTProviders": [
              {
                  "iss": "http://issuer.com/connect/token",
                  "aud": "my-service-audience"
              }
            ],
            "UnauthenticatedPaths": [
                "/health"
            ]
        }
---

```

## Roadmap

- Add metrics
- Improve/add new configuration methods

## Security

We take security seriously at Soluto. In case you find a security issue or have something you would like to discuss refer to our [Security.md](SECURITY) policy.

## Contributing

If you found a bug or have idea for cool feature please open issue and let us know.
> Please notice: Do not report security issues on GitHub. We will immediately delete such issues.

## License

Licensed under [the MIT License](LICENSE)
