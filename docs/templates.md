Kontemplate templates
=====================

The template file format is based on Go's [templating engine][] in combination
with a small extension library called [sprig][] that adds additional template
functions.

Go templates can either simply display variables or build more complicated
*pipelines* in which variables are passed to functions for further processing,
or in which conditionals are evaluated for more complex template logic.

It is recommended that you check out the Golang [documentation][] for the templating
engine in addition to the cherry-picked features listed here.

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [Kontemplate templates](#kontemplate-templates)
    - [Basic variable interpolation](#basic-variable-interpolation)
        - [Example:](#example)
    - [Template functions](#template-functions)
    - [Examples:](#examples)
    - [Conditionals & ranges](#conditionals--ranges)
    - [Caveats](#caveats)

<!-- markdown-toc end -->

## Basic variable interpolation

The basic template format uses `{{ .variableName }}` as the interpolation format.

### Example:

Assuming that you include a resource set as such:

```
- name: api-gateway
  values:
    internalHost: http://my-internal-host/
```

And the api-gateway resource set includes a ConfigMap (some fields left out for
the example):

```
# api-gateway/configmap.yaml:
---
kind: ConfigMap
metadata:
  name: api-gateway-config
data:
  internalHost: {{ .internalHost }}
```

The resulting output will be:

```

---
kind: ConfigMap
metadata:
  name: api-gateway-config
data:
  internalHost: http://my-internal-host/
```

## Template functions

Go templates support template functions which you can think of as a sort of
shell-like pipeline where text flows through transformations from left to
right.

Some template functions come from Go's standard library and are listed in the
[Go documentation][]. In addition the functions declared by [sprig][] are
available in kontemplate, as well as two custom functions:

`json`: Encodes any supplied data structure as JSON.
`passLookup`: Looks up the supplied key in [pass][]

## Examples:

```
# With the following values:
name: Donald
certKeyPath: my-website/cert-key

# The following interpolations are possible:

{{ .name | upper }}
-> DONALD

{{ .name | upper | repeat 2 }}
-> DONALD DONALD

{{ .certKeyPath | passLookup }}
-> Returns content of 'my-website/cert-key' from pass
```

## Conditionals & ranges

Some logic is supported in Golang templates and can be used in Kontemplate, too.

With the following values:

```
useKube2IAM: true
servicePorts:
  - 8080
  - 9090
```

The following interpolations are possible:

```
# Conditionally insert something in the template:
metadata:
  annotations:
    foo: bar
    {{ if .useKube2IAM -}} iam.amazonaws.com/role: my-api {{- end }}
```

```
# Iterate over a list of values
ports:
  {{ range .servicePorts }}
  - port: {{ . }}
  {{ end }}
```

Check out the Golang documentation (linked above) for more information about template logic.

## Caveats

Kontemplate does not by itself parse any of the content of the templates, which
means that it does not validate whether the resources you supply are valid YAML
or JSON.

You can perform some validation by using `kontemplate apply --dry-run` which
will make use of the Dry-Run functionality in `kubectl`.

[templating engine]: https://golang.org/pkg/text/template/
[documentation]: https://golang.org/pkg/text/template/
[sprig]: http://masterminds.github.io/sprig/
[Go documentation]: https://golang.org/pkg/text/template/#hdr-Functions
[pass]: https://www.passwordstore.org/
