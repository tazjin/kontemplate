Kontemplate tips & tricks
=========================

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [Kontemplate tips & tricks](#kontemplate-tips--tricks)
    - [Update Deployments when ConfigMaps change](#update-deployments-when-configmaps-change)
    - [direnv & pass](#direnv--pass)

<!-- markdown-toc end -->

## Update Deployments when ConfigMaps change

Kubernetes does [not currently][] have the ability to perform rolling updates
of Deployments and other resource types when `ConfigMap` or `Secret` objects
are updated.

It is possible to make use of annotations and templating functions in
Kontemplate to force updates to these resources anyways.
 
For example:

```yaml
# A ConfigMap that contains some configuration for your app
---
kind: ConfigMap
metadata:
  name: app-config
data:
  app.conf: |
    name: {{ .appName }}
    foo: bar
```

Now whenever the `appName` variable changes or we make an edit to the
`ConfigMap` we would like to update the `Deployment` making use of it, too. We
can do this by adding a hash of the parsed template to the annotations of the
created `Pod` objects:

```yaml

---
kind: Deployment
metadata:
  name: app
spec:
  template:
    metadata:
      annotations:
        configHash: {{ insertTemplate "app-config.yaml" | sha256sum }}
    spec:
      containers:
        - name: app
          # Some details omitted ... 
          volumeMounts:
            - name: config
              mountPath: /etc/app/
      volumes:
        - name: config
          configMap:
            name: app-config
```

Now any change to the `ConfigMap` - either by directly editing the yaml file or
via a changed template variable - will cause the annotation to change,
triggering a rolling update of all relevant pods.

## direnv & pass

Users of `pass` may have multiple different password stores on their machines.
Assuming that `kontemplate` configuration exists somewhere on the filesystem
per project, it is easy to use [direnv][] to switch to the correct
`PASSWORD_STORE_DIR` variable when entering the folder.

[not currently]: https://github.com/kubernetes/kubernetes/issues/22368
[direnv]: https://direnv.net/
