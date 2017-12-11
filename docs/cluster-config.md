Cluster configuration
==========================

Every cluster (or "environment") that requires individual configuration is specified in
a very simple YAML file in Kontemplate.

An example file for a hypothetical test environment could look like this:

```yaml
---
context: k8s.test.mydomain.com
global:
  clusterName: test-cluster
  defaultReplicas: 2
import:
  - test-secrets.yaml
include:
  - name: gateway
    path: tools/nginx
    values:
      tlsDomains:
        - test.oslo.pub
        - test.tazj.in
  - path: backend
    values:
      env: test
    include:
      - name: blog
        values:
          url: test.tazj.in
      - name: pub-service
```

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [Cluster configuration](#cluster-configuration)
    - [Fields](#fields)
        - [`context`](#context)
        - [`global`](#global)
        - [`import`](#import)
        - [`include`](#include)
    - [External variables](#external-variables)

<!-- markdown-toc end -->

## Fields

This is documentation for the individual fields in a cluster context file.

### `context`

The `context` field contains the name of the kubectl-context. You can list context names with
'kubectl config get-contexts'.

This must be set here so that Kontemplate can use the correct context when calling kubectl.

This field is **required** for `kubectl`-wrapping commands. It can be left out if only the `template`-command is used.

### `global`

The `global` field contains a key/value map of variables that should be available to all resource
sets in the cluster.

This field is **optional**.

### `import`

The `import` field contains the file names of additional YAML or JSON files from which global
variables should be loaded. Using this field makes it possible to keep certain configuration that
is the same for some, but not all, clusters in a common place.

This field is **optional**.

### `include`

The `include` field contains the actual resource sets to be included in the cluster.

Information about the structure of resource sets can be found in the [resource set documentation][].

This field is **required**.

## External variables

As mentioned above, extra variables can be loaded from additional YAML or JSON files. Assuming you
have a file called `test-secrets.yaml` which contains variables that should be shared between a `test`
and `dev` cluster, you could include it in your context as such:

```yaml
# test-secrets.yaml:
mySecretVar: foo-bar-12345

# test-cluster.yaml:
context: k8s.test.mydomain.com
include:
  - test-secrets.yaml

# dev-cluster.yaml:
context: k8s.dev.mydomain.com
include:
  - test-secrets.yaml
```

The variable `mySecretVar` is then available as a global variable.

[resource set documentation]: resource-sets.md
