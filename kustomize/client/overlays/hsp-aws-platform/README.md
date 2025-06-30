# HSP AWS Platform Overlay

This kustomize overlay configures the dip-forwarder-client for deployment on the HSP AWS Platform.

## Features

- Adds the `-use-private-link` command argument to enable private link connectivity
- Applies the `hsp-aws-platform-` prefix to resource names for environment isolation

## Usage

To build and apply this overlay:

```bash
# Build the configuration
kustomize build kustomize/client/overlays/hsp-aws-platform/

# Apply to Kubernetes cluster
kustomize build kustomize/client/overlays/hsp-aws-platform/ | kubectl apply -f -
```

## Configuration

The overlay modifies the base deployment by:

- Adding the `-use-private-link` argument to the container command

This allows the client to connect through AWS private links when deployed in the HSP AWS environment.
