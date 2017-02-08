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

```
NAME:
   kontemplate - simple Kubernetes resource templating

USAGE:
   kontemplate [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     template  Interpolate and print templates
     apply     Interpolate templates and run 'kubectl apply'
     replace   Interpolate templates and run 'kubectl replace'
     delete    Interpolate templates and run 'kubectl delete'
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

All options support the same set of extra flags:

```
OPTIONS:
   --file value, -f value     Cluster configuration file to use
   --include value, -i value  Limit templating to explicitly included resource sets
   --exclude value, -e value  Exclude certain resource sets from templating
```

Examples:

```
# Look at output for a specific resource set and check to see if it's correct ...
kontemplate template -f example/prod-cluster.yaml -i some-api

# ... maybe do a dry-run to see what kubectl would do:
kontemplate apply -f example/prod-cluster.yaml --dry-run

# And actually apply it if you like what you see:
kontemplate apply -f example/prod-cluster.yaml
```