# Nango Integration Operator

A Kubernetes operator for managing Nango integrations declaratively. This operator allows you to create and manage Nango integrations using Kubernetes Custom Resources.

## Overview

The Nango Integration Operator provides a Kubernetes-native way to manage Nango integrations. It watches for `NangoIntegration` custom resources and automatically creates the corresponding integrations in your Nango instance.

## Features

- **Declarative Integration Management**: Define integrations as Kubernetes resources
- **Automatic Reconciliation**: The operator continuously ensures the desired state
- **Status Tracking**: Monitor integration creation status and health
- **Error Handling**: Automatic retry logic for failed operations

## Installation

### Prerequisites

- Kubernetes cluster (1.19+)
- kubectl configured to communicate with your cluster
- Nango instance with API access

### Deploy the Operator

1. Clone this repository:
```bash
git clone <repository-url>
cd nango-integration-operator
```

2. Deploy the operator:
```bash
make deploy
```

3. Verify the deployment:
```bash
kubectl get pods -n nango-integration-operator-system
```

## Usage

### Creating a Nango Integration

Create a `NangoIntegration` resource with the required configuration:

```yaml
apiVersion: nango.nango.dev/v1alpha1
kind: NangoIntegration
metadata:
  name: my-slack-integration
  namespace: default
spec:
  unique_key: "slack-my-app"
  provider: "slack"
  display_name: "Slack Integration"
  credentials:
    type: "OAUTH1"
    client_id: "your-slack-client-id"
    client_secret: "your-slack-client-secret"
    scopes: "chat:write,channels:read"
  nango_token: "your-nango-api-token"
  nango_base_url: "https://api.nango.dev"  # Optional, defaults to this value
```

### Required Fields

- `unique_key`: A unique identifier for the integration
- `provider`: The provider name (e.g., "slack", "github", "salesforce")
- `display_name`: Human-readable name for the integration
- `credentials`: OAuth credentials for the integration
  - `type`: Authentication type (e.g., "OAUTH1", "OAUTH2")
  - `client_id`: OAuth client ID (can be a string value or secretKeyRef)
  - `client_secret`: OAuth client secret (can be a string value or secretKeyRef)
  - `scopes`: Required OAuth scopes (optional)

### Optional Fields

- `nango_token`: Your Nango API token for authentication (can be a string value or secretKeyRef)
- `nango_base_url`: Custom Nango API base URL (defaults to https://api.nango.dev)

### Checking Integration Status

Monitor the status of your integration:

```bash
kubectl get nangointegration my-slack-integration -o yaml
```

The status will show:
- `status.status`: Current status ("Created", "Failed", etc.)
- `status.error_message`: Error details if creation failed
- `status.conditions`: Kubernetes conditions for monitoring

### Example: Slack Integration

```yaml
apiVersion: nango.nango.dev/v1alpha1
kind: NangoIntegration
metadata:
  name: slack-integration
spec:
  unique_key: "slack-nango-community"
  provider: "slack"
  display_name: "Slack"
  credentials:
    type: "OAUTH1"
    client_id: "123456789.123456789"
    client_secret: "your-slack-client-secret"
    scopes: "chat:write,channels:read,users:read"
  nango_token: "nango_sk_..."
```

### Example: GitHub Integration

```yaml
apiVersion: nango.nango.dev/v1alpha1
kind: NangoIntegration
metadata:
  name: github-integration
spec:
  unique_key: "github-my-app"
  provider: "github"
  display_name: "GitHub"
  credentials:
    type: "OAUTH2"
    client_id: "your-github-client-id"
    client_secret: "your-github-client-secret"
    scopes: "repo,user"
  nango_token: "nango_sk_..."
```

## Security Considerations

### Secret References (Recommended for Production)

The operator supports referencing sensitive values from Kubernetes secrets. This is the recommended approach for production deployments.

#### Creating Secrets

First, create the secrets with your credentials:

```yaml
# OAuth credentials secret
apiVersion: v1
kind: Secret
metadata:
  name: slack-oauth-secret
  namespace: default
type: Opaque
data:
  client_id: <base64-encoded-client-id>
  client_secret: <base64-encoded-client-secret>

---
# Nango API token secret
apiVersion: v1
kind: Secret
metadata:
  name: nango-api-secret
  namespace: default
type: Opaque
data:
  token: <base64-encoded-nango-token>
```

#### Using Secret References

Reference the secrets in your NangoIntegration:

```yaml
apiVersion: nango.nango.dev/v1alpha1
kind: NangoIntegration
metadata:
  name: secure-slack-integration
spec:
  unique_key: "slack-secure"
  provider: "slack"
  display_name: "Slack Integration"
  credentials:
    type: "OAUTH1"
    client_id:
      secretKeyRef:
        name: "slack-oauth-secret"
        key: "client_id"
    client_secret:
      secretKeyRef:
        name: "slack-oauth-secret"
        key: "client_secret"
    scopes: "chat:write,channels:read"
  nango_token:
    secretKeyRef:
      name: "nango-api-secret"
      key: "token"
```

#### Using Direct Values (for development/testing)

```yaml
apiVersion: nango.nango.dev/v1alpha1
kind: NangoIntegration
metadata:
  name: dev-slack-integration
spec:
  unique_key: "slack-dev"
  provider: "slack"
  display_name: "Slack Integration (Dev)"
  credentials:
    type: "OAUTH1"
    client_id:
      value: "your-client-id"
    client_secret:
      value: "your-client-secret"
    scopes: "chat:write,channels:read"
  nango_token:
    value: "your-nango-token"
```

#### Encoding Values

To encode your values for the secret:

```bash
echo -n "your-client-id" | base64
echo -n "your-client-secret" | base64
echo -n "your-nango-token" | base64
```

## Troubleshooting

### Check Operator Logs

```bash
kubectl logs -n nango-integration-operator-system deployment/nango-integration-operator-controller-manager
```

### Common Issues

1. **Authentication Errors**: Ensure your Nango API token is valid
2. **Invalid Credentials**: Verify OAuth client ID and secret are correct
3. **Network Issues**: Check connectivity to the Nango API endpoint

### Status Conditions

The operator sets the following conditions:
- `Ready`: Indicates if the integration is successfully created
- `IntegrationCreated`: Success condition
- `IntegrationCreationFailed`: Failure condition

## Development

### Building from Source

```bash
make build
```

### Running Tests

```bash
make test
```

### Running Locally

```bash
make run
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.

