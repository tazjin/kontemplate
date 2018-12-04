Resource Sets
================

Resource sets are collections of Kubernetes resources that should be passed to `kubectl` together.

Technically a resource set is simply a folder with a few YAML and/or JSON templates in it.

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [Resource Sets](#resource-sets)
- [Creating resource sets](#creating-resource-sets)
    - [Default variables](#default-variables)
- [Including resource sets](#including-resource-sets)
    - [Fields](#fields)
        - [`name`](#name)
        - [`path`](#path)
        - [`values`](#values)
        - [`include`](#include)
        - [`helper`](#helper)
    - [Multiple includes](#multiple-includes)
    - [Nesting resource sets](#nesting-resource-sets)
        - [Caveats](#caveats)

<!-- markdown-toc end -->

# Creating resource sets

Simply create a folder in your Kontemplate repository and place a YAML or JSON file in it. These
files get interpreted as [templates][] during Kontemplate runs and variables (as well as template
logic or functions) will be interpolated.

Refer to the template documentation for information on how to write templates.

## Default variables

Sometimes it is useful to specify default values for variables that should be interpolated during
a run if the [cluster configuration][] does not specify a variable explicitly.

This can be done simply by placing a `default.yaml` or `default.json` file in the resource set
folder and filling it with key/value pairs of the intended default variables.

Kontemplate will error during interpolation if any variables are left unspecified.

# Including resource sets

Under the cluster configuration `include` key resource sets are included and required variables
are specified. For example:

```yaml
include:
  - name: some-api
    values:
      version: 1.2-SNAPSHOT
```

This will include a resource set from a folder called `some-api` and set the specified `version` variable.

## Fields

The available fields when including a resource set are these:

### `name`

The `name` field contains the name of the resource set. This name can be used to refer to the resource set
when specifying explicit includes or excludes during a run.

By default it is assumed that the `name` is the path to the resource set folder, but this can be overridden.

This field is **required**.

### `path`

The `path` field specifies an explicit path to a resource set folder in the case that it should differ from
the resource set's `name`.

This field is **optional**.

### `values`

The `values` field specifies key/values pairs of variables that should be available during templating.

This field is **optional**.

### `include`

The `include` field specifies additional resource sets that should be included and that should inherit the
variables of this resource set.

The fully qualified names of "nested" resource sets are set to `${PARENT_NAME}/${CHILD_NAME}` and paths are
merged in the same way.

This makes it easy to organise different resource sets as "groups" to include / exclude them collectively
during runs.

This field is **optional**.

### `helper`

The `helper` field specifies additional templates to load when rendering the final ResourceSet template
file.

This fields allow you to `define` named golang template and include them with `template "templateName"`
function.

This field is **optional**.

## Multiple includes

Resource sets can be included multiple times with different configurations. In this case it is recommended
to set the `path` and `name` fields explicitly. For example:

```yaml
include:
  - name: forwarder-europe
    path: tools/forwarder
    values:
      source: europe
  - name: forwarder-asia
    path: tools/forwarder
    values:
      source: asia
```

The two different configurations can be referred to by their set names, but will use the same resource
templates with different configurations.

## Nesting resource sets

As mentioned above for the `include` field, resource sets can be nested. This lets users group resource
sets in logical ways using simple folder structures.

Assuming a folder structure like:

```
├── backend
│   ├── auth-api
│   ├── message-api
│   └── order-api
└── frontend
    ├── app-page
    └── login-page
```

With each of these folders being a resource set, they could be included in a cluster configuration like so:

```yaml
include:
  - name: backend
    include:
      - name: auth-api
      - name: message-api
      - name: order-api
  - name: frontend:
    include:
      - name: app-page
      - name: login-page
```

Kontemplate could then be run with, for example, `--include backend` to only include the resource sets nested
in the backend group. Specific resource sets can also be targeted, for example as `--include backend/order-api`.

Variables specified in the parent resource set are inherited by the children.

### Caveats

Two caveats apply that users should be aware of:

1. The parent resource set can not contain any resource templates itself.

2. Only one level of nesting is supported. Specifying `include` again on a nested resource set will be ignored.

[templates]: templates.md
[cluster configuration]: cluster-config.md
