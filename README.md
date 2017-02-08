KonTemplate - A simple Kubernetes templater
===========================================

[![Build Status](https://travis-ci.org/tazjin/kontemplate.svg?branch=master)](https://travis-ci.org/tazjin/kontemplate)

I made this tool out of frustration with the available ways to template Kubernetes resource files. All I want out of
such a tool is a way to specify lots of resources with placeholders that get filled in with specific values, based on
which context (i.e. k8s cluster) is specified.

## Overview

KonTemplate lets you describe resources as you normally would in a simple folder structure:

```
.
├── prod-cluster.yaml
└── some-api
    ├── deployment.yaml
    └── service.yaml
```

This example has all resources belonging to `some-api` (no file naming conventions enforced at all!) in the `some-api`
folder and the configuration for the cluster `prod-cluster` in the corresponding file.

Lets take a short look at `prod-cluster.yaml`:

```yaml
---
context: k8s.prod.mydomain.com
global:
  globalVar: lizards
include:
  - name: some-api
    values:
      version: 1.0-0e6884d
      importantFeature: true
      apiPort: 4567
```

Those values are then templated into the resource files of `some-api`.

## Installation

Assuming you have Go configured correctly, you can simply `go get github.com/tazjin/kontemplate/...`.

## Usage

You must have `kubectl` installed to use KonTemplate effectively.

At the moment KonTemplate will simply output the templated Kubernetes resource files, which can
then be piped into `kubectl`:

```
# Look at output and check to see if it's correct ...
kontemplate run -f example/prod-cluster.yaml -l some-api

# ... if it is, go ahead and apply it
kontemplate run -f example/prod-cluster.yaml -l some-api | kubectl apply -f -

# That's it!
```