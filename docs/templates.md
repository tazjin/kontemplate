Kontemplate templates
=====================

The template file format is based on Go's [templating engine][] in combination
with a small extension library called [sprig][] that adds additional template
functions.

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

## Caveats

Kontemplate does not by itself parse any of the content of the templates, which
means that it does not validate whether the resources you supply are valid YAML
or JSON.

You can perform some validation by using `kontemplate apply --dry-run` which
will make use of the Dry-Run functionality in `kubectl`.

[templating engine]: https://golang.org/pkg/text/template/
[sprig]: http://masterminds.github.io/sprig/
[Go documentation]: https://golang.org/pkg/text/template/#hdr-Functions
[pass]: https://www.passwordstore.org/
