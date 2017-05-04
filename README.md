Kontemplate - A simple Kubernetes templater
===========================================

[![Build Status](https://travis-ci.org/tazjin/kontemplate.svg?branch=master)](https://travis-ci.org/tazjin/kontemplate)

I made this tool out of frustration with the available ways to template Kubernetes resource files. All I want out of
such a tool is a way to specify lots of resources with placeholders that get filled in with specific values, based on
which context (i.e. k8s cluster) is specified.

## Overview

Kontemplate lets you describe resources as you normally would in a simple folder structure:

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

You must have `kubectl` installed to use Kontemplate effectively.

```
usage: kontemplate [<flags>] <command> [<args> ...]

simple Kubernetes resource templating

Flags:
  -h, --help                 Show context-sensitive help (also try --help-long and --help-man).
  -i, --include=INCLUDE ...  Resource sets to include explicitly
  -e, --exclude=EXCLUDE ...  Resource sets to exclude explicitly

Commands:
  help [<command>...]
    Show help.

  template <file>
    Template resource sets and print them

  apply [<flags>] <file>
    Template resources and pass to 'kubectl apply'

  replace <file>
    Template resources and pass to 'kubectl replace'

  delete <file>
    Template resources and pass to 'kubectl delete'

  create <file>
    Template resources and pass to 'kubectl create'

```

Examples:

```
# Look at output for a specific resource set and check to see if it's correct ...
kontemplate template example/prod-cluster.yaml -i some-api

# ... maybe do a dry-run to see what kubectl would do:
kontemplate apply example/prod-cluster.yaml --dry-run

# And actually apply it if you like what you see:
kontemplate apply example/prod-cluster.yaml
```