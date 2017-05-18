Kontemplate - A simple Kubernetes templater
===========================================

[![Build Status](https://travis-ci.org/tazjin/kontemplate.svg?branch=master)](https://travis-ci.org/tazjin/kontemplate)

Kontemplate is a simple CLI tool that can take sets of Kubernetes resource
files with placeholders and insert values per environment.

This tool was made because in many cases all I want in terms of Kubernetes
configuration is simple value interpolation per environment (i.e. Kubernetes
cluster), but with the same deployment files.

In my experience this is often enough and more complex solutions such as
[Helm][] are not required.

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

Those values are then templated into the resource files of `some-api`. That's it!

You can also set up more complicated folder structures for organisation, for example:

```
.
├── api
│   ├── image-api
│   │   └── deployment.yaml
│   └── music-api
│       └── deployment.yaml
│   │   └── default.json
├── frontend
│   ├── main-app
│   │   ├── deployment.yaml
│   │   └── service.yaml
│   └── user-page
│       ├── deployment.yaml
│       └── service.yaml
├── prod-cluster.yaml
└── test-cluster.yaml
```

And selectively template or apply resources with a command such as
`kontemplate apply test-cluster.yaml --include api --include frontend/user-page`
to only update the `api` resource sets and the `frontend/user-page` resource set.

## Installation

It is recommended to install Kontemplate from the signed binary releases available on the
[releases page][]. Release binaries are available for Linux, OS X, FreeBSD and Windows.

### Homebrew

OS X users with Homebrew installed can "tap" Kontemplate like such:

```sh
brew tap tazjin/kontemplate https://github.com/tazjin/kontemplate
brew install kontemplate
```

### Arch Linux

An [AUR package][] is available for Arch Linux and other `pacman`-based distributions.

### Building repeatably from source

Version pinning for Go dependencies is provided by a [Repeatr][] formula. After cloning
the repository the latest release can be built with `repeatr run kontemplate.frm`.

This will place release binaries in the `release` folder.

### Building from source

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

[Helm]: https://helm.sh/
[releases page]: https://github.com/tazjin/kontemplate/releases
[AUR package]: https://aur.archlinux.org/packages/kontemplate-git/
[Repeatr]: http://repeatr.io/
