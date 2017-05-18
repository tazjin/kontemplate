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
    hash: "1487840a4bf30270decdc04123c41cdc7a8067c9"
    silo: "https://github.com/polydawn/meep"
  "/go/src/github.com/Masterminds/sprig":
    type: "git"
    hash: "f5b0ed4a680a0943228155eaf6a77a96ead1bc77"
    silo: "https://github.com/Masterminds/sprig"
  "/go/src/github.com/ghodss/yaml":
    type: "git"
    hash: "0ca9ea5df5451ffdf184b4428c902747c2c11cd7"
    silo: "https://github.com/ghodss/yaml"
  "/go/src/gopkg.in/yaml.v2":
    type: "git"
    hash: "cd8b52f8269e0feb286dfeef29f8fe4d5b397e0b"
    silo: "https://gopkg.in/yaml.v2"
  "/go/src/gopkg.in/alecthomas/kingpin.v2":
    type: "git"
    hash: "7f0871f2e17818990e4eed73f9b5c2f429501228"
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
    hash: "abff1900528dbdaf6f3f5aa92c398be1eaf2a9f7"
    silo: "https://github.com/Masterminds/semver"
  "/go/src/github.com/aokoli/goutils":
    type: "git"
    hash: "e57d01ace047c1a43e6a49ecf3ecc50ed2be81d1"
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
      git clone --single-branch --branch v1.0.2 https://github.com/tazjin/kontemplate /go/src/github.com/tazjin/kontemplate
      cd /go/src/github.com/tazjin/kontemplate
      ./build-release.sh build
outputs:
  "release":
    type: "dir"
    mount: "/go/src/github.com/tazjin/kontemplate/release"
    silo: "file://release"
