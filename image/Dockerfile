FROM alpine:3.10

ADD hashes /root/hashes
ADD https://storage.googleapis.com/kubernetes-release/release/v1.15.3/bin/linux/amd64/kubectl /usr/bin/kubectl
ADD https://github.com/tazjin/kontemplate/releases/download/v1.7.0/kontemplate-1.7.0-511ae92-linux-amd64.tar.gz /tmp/kontemplate.tar.gz

# Pass release version is 1.7.3
ADD https://raw.githubusercontent.com/zx2c4/password-store/74fdfb5022f317ad48d449e29543710bdad1afda/src/password-store.sh /usr/bin/pass

RUN sha256sum -c /root/hashes && \
    apk add -U bash tree gnupg git && \
    chmod +x /usr/bin/kubectl /usr/bin/pass && \
    tar xzvf /tmp/kontemplate.tar.gz && \
    mv kontemplate /usr/bin/kontemplate && \
    /usr/bin/kontemplate version
