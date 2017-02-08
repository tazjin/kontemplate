KonTemplate - A simple Kubernetes templater
===========================================

I made this tool out of frustration with the available ways to template Kubernetes resource files. All I want out of
such a tool is a way to specify lots of resources with placeholders that get filled in with specific values, based on
which context (i.e. k8s cluster) is specified.

## Overview

KonTemplate lets you describe resources as you normally would in a simple folder structure:

```
.
├── prod-cluster.json
└── some-api
    ├── deployment.yaml
    └── service.yaml
```

This example has all resources belonging to `some-api` (no file naming conventions enforced at all!) in the `some-api`
folder and the configuration for the cluster `prod-cluster` in the corresponding file.

Lets take a short look at `prod-cluster.json`:

```json
{
  "context": "k8s.prod.mydomain.com",
  "include": [
    {
      "name": "some-api",
      "values": {
        "importantFeature": true,
        "apiPort": 4567
      }
    }
  ]
}
```


Those values are then templated into the resource files of `some-api`.

## Usage

You must have `kubectl` installed to use KonTemplate.
