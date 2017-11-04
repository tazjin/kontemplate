inputs:
  "/": # https://github.com/tklx/base
    type: "tar"
    hash: "9nkvYhmJHaeK_Agc3Lm5rg444dSLWDp0Pri-KilHiX3A9Pt4TaQ7RxOj5qMSs6XT"
    silo: "https://github.com/tklx/base/releases/download/0.1.1/rootfs.tar.xz"
  "/opt":
    type: "tar"
    hash: "gi0Kpb-VH3TK0UBX6YmpuKsrMAUlxicPrY2YvXPo9sBQm_NsD_hKrn7pmc95zrmM"
    silo: "https://storage.googleapis.com/golang/go1.8.1.linux-amd64.tar.gz"
  # Kontemplate dependencies!
  "/go/src/github.com/polydawn/meep":
    type: "git"
    hash: "eaf1db2168fe380b4da17a35f0adddb5ae15a651"
    silo: "https://github.com/polydawn/meep"
  "/go/src/github.com/Masterminds/sprig":
    type: "git"
    hash: "e039e20e500c2c025d9145be375e27cf42a94174"
    silo: "https://github.com/Masterminds/sprig"
  "/go/src/github.com/ghodss/yaml":
    type: "git"
    hash: "0ca9ea5df5451ffdf184b4428c902747c2c11cd7"
    silo: "https://github.com/ghodss/yaml"
  "/go/src/gopkg.in/yaml.v2":
    type: "git"
    hash: "eb3733d160e74a9c7e442f435eb3bea458e1d19f"
    silo: "https://gopkg.in/yaml.v2"
  "/go/src/gopkg.in/alecthomas/kingpin.v2":
    type: "git"
    hash: "1087e65c9441605df944fb12c33f0fe7072d18ca"
    silo: "https://gopkg.in/alecthomas/kingpin.v2"
  "/go/src/github.com/alecthomas/template":
    type: "git"
    hash: "a0175ee3bccc567396460bf5acd36800cb10c49c"
    silo: "https://github.com/alecthomas/template"
  "/go/src/github.com/alecthomas/units":
    type: "git"
    hash: "2efee857e7cfd4f3d0138cc3cbb1b4966962b93a"
    silo: "https://github.com/alecthomas/units"
  "/go/src/github.com/Masterminds/semver":
    type: "git"
    hash: "517734cc7d6470c0d07130e40fd40bdeb9bcd3fd"
    silo: "https://github.com/Masterminds/semver"
  "/go/src/github.com/aokoli/goutils":
    type: "git"
    hash: "3391d3790d23d03408670993e957e8f408993c34"
    silo: "https://github.com/aokoli/goutils"
  "/go/src/github.com/huandu/xstrings":
    type: "git"
    hash: "3959339b333561bf62a38b424fd41517c2c90f40"
    silo: "https://github.com/huandu/xstrings"
  "/go/src/github.com/imdario/mergo":
    type: "git"
    hash: "d806ba8c21777d504a2090a2ca4913c750dd3a33"
    silo: "https://github.com/imdario/mergo"
  "/go/src/github.com/satori/go.uuid":
    type: "git"
    hash: "5bf94b69c6b68ee1b541973bb8e1144db23a194b"
    silo: "https://github.com/satori/go.uuid"
  "/go/src/golang.org/x/crypto":
    type: "git"
    hash: "ab89591268e0c8b748cbe4047b00197516011af5"
    silo: "https://go.googlesource.com/crypto"
action:
  policy: governor
  command:
    - "/bin/sh"
    - "-e"
    - "-c"
    - |
      export PATH="/opt/go/bin:$PATH"
      export GOROOT=/opt/go
      export GOPATH=/go
      echo 'nameserver 8.8.8.8' > /etc/resolv.conf
      apt-get update && apt-get install -y git ca-certificates
      mkdir -p /go/src/github.com/tazjin
      git clone --single-branch --branch v1.3.0 https://github.com/tazjin/kontemplate /go/src/github.com/tazjin/kontemplate
      cd /go/src/github.com/tazjin/kontemplate
      ./build-release.sh build
outputs:
  "release":
    type: "dir"
    mount: "/go/src/github.com/tazjin/kontemplate/release"
    silo: "file://release"
