Kontemplate Docker image
========================

This builds a simple Docker image available on the Docker Hub as `tazjin/kontemplate`.

Builds are automated based on the Dockerfile contained here.

It contains both `kontemplate` and `kubectl` and can be used as part of container-based
CI pipelines.

`pass` and its dependencies are also installed to enable the use of the `passLookup`
template function if desired.
