Kontemplate tips & tricks
=========================


## Update Deployments when ConfigMaps change

Kubernetes does [not currently][] have the ability to perform rolling updates
of Deployments and other resource types when `ConfigMap` or `Secret` objects
are updated.

It is possible to make use of annotations and templating functions in
Kontemplate to force updates to these resources anyways (assuming that the
`ConfigMap` or `Secret` contains interpolated variables).
 
For example:

```yaml
# A ConfigMap that contains some data structure in JSON format
---
kind: ConfigMap
metadata:
  name: app-config
data:
  configFile: {{ .appConfig | json }}
```

Now whenever the `appConfig` variable changes we would like to update the
`Deployment` making use of it, too. We can do this by adding a hash of the
configuration to the annotations of the created `Pod` objects:

```yaml

---
kind: Deployment
metadata:
  name: app
spec:
  template:
    metadata:
      annotations:
        configHash: {{ .appConfig | json | sha256sum }}
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

Now if the `ConfigMap` object appears first in the resource files, `kubectl`
will apply the resources sequentially and the updated annotation will cause
a rolling update of all relevant pods.

## direnv & pass

Users of `pass` may have multiple different password stores on their machines.
Assuming that `kontemplate` configuration exists somewhere on the filesystem
per project, it is easy to use [direnv][] to switch to the correct
`PASSWORD_STORE_DIR` variable when entering the folder.

[not currently]: https://github.com/kubernetes/kubernetes/issues/22368
[direnv]: https://direnv.net/