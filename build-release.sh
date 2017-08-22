#!/bin/bash
set -ueo pipefail

readonly GIT_HASH="$(git rev-parse --short HEAD)"
readonly LDFLAGS="-X main.gitHash=${GIT_HASH} -w -s"
readonly VERSION="1.1.0-${GIT_HASH}"

function binary-name() {
    local os="${1}"
    local target="${2}"
    if [ "${os}" = "windows" ]; then
        echo -n "${target}/kontemplate.exe"
    else
        echo -n "${target}/kontemplate"
    fi
}

function build-for() {
    local os="${1}"
    local arch="${2}"
    local target="release/${os}/${arch}"
    local bin=$(binary-name "${os}" "${target}")

    echo "Building kontemplate for ${os}-${arch} in ${target}"

    mkdir -p "${target}"

    env GOOS="${os}" GOARCH="${arch}" go build \
        -ldflags "${LDFLAGS}" \
        -o "${bin}" \
        -tags netgo
}

function sign-for() {
    local os="${1}"
    local arch="${2}"
    local target="release/${os}/${arch}"
    local bin=$(binary-name "${os}" "${target}")
    local tar="release/kontemplate-${VERSION}-${os}-${arch}.tar.gz"

    echo "Packing release into ${tar}"
    tar czvf "${tar}" -C "${target}" $(basename "${bin}")

    local hash=$(sha256sum "${tar}")
    echo "Signing kontemplate release tarball for ${os}-${arch} with SHA256 ${hash}"
    gpg --armor --detach-sig --sign "${tar}"
}

case "${1}" in
    "build")
        # Build releases for various operating systems:
        build-for "linux" "amd64"
        build-for "darwin" "amd64"
        build-for "windows" "amd64"
        build-for "freebsd" "amd64"
        exit 0
        ;;
    "sign")
        # Bundle and sign releases:
        sign-for "linux" "amd64"
        sign-for "darwin" "amd64"
        sign-for "windows" "amd64"
        sign-for "freebsd" "amd64"
        exit 0
        ;;
esac
