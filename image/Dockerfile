FROM alpine:3.8

ADD hashes /root/hashes
ADD https://storage.googleapis.com/kubernetes-release/release/v1.13.2/bin/linux/amd64/kubectl /usr/bin/kubectl
ADD https://github.com/tazjin/kontemplate/releases/download/v1.7.0/kontemplate-1.7.0-511ae92-linux-amd64.tar.gz /tmp/kontemplate.tar.gz

# Pass release version is 1.7.1
ADD https://raw.githubusercontent.com/zx2c4/password-store/38ec1c72e29c872ec0cdde82f75490640d4019bf/src/password-store.sh /usr/bin/pass

RUN sha256sum -c /root/hashes && \
    apk add -U bash tree gnupg git && \
    chmod +x /usr/bin/kubectl /usr/bin/pass && \
    tar xzvf /tmp/kontemplate.tar.gz && \
    mv kontemplate /usr/bin/kontemplate && \
    /usr/bin/kontemplate version
