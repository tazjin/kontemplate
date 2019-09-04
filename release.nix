# Copyright (C) 2016-2019  Vincent Ambo <mail@tazj.in>
#
# This file is part of Kontemplate.
#
# Kontemplate is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This file is the Nix derivation used to build release binaries for
# several different architectures and operating systems.

let pkgs = import ((import <nixpkgs> {}).fetchFromGitHub {
  owner  = "NixOS";
  repo   = "nixpkgs";
  rev    = "50f41ea2fcf86def32799f75577a4fe5cfd1132e";
  sha256 = "1q0bxl5nxx1kabqvyzkdw91c5dnwpi2rwsgs5jdmnj7f0qqgdxh8";
}) {};
in with pkgs; buildGoPackage rec {
  name = "kontemplate-${version}";
  version = "master";
  src = ./.;
  goPackagePath = "github.com/tazjin/kontemplate";
  goDeps = ./deps.nix;

  # This configuration enables the building of statically linked
  # executables. For some reason, those will have multiple references
  # to the Go compiler's installation path in them, which is the
  # reason for setting the 'allowGoReference' flag.
  dontStrip = true; # Linker configuration handles stripping
  allowGoReference = true;
  CGO_ENABLED="0";
  GOCACHE="off";

  # Configure release builds via the "build-matrix" script:
  buildInputs = [ git ];
  buildPhase = ''
    cd go/src/${goPackagePath}
    ./build-release.sh build
  '';

  outputs = [ "out" ];
  installPhase = ''
    mkdir $out
    cp -r release/ $out
  '';

  meta = with lib; {
    description = "A resource templating helper for Kubernetes";
    homepage = "http://kontemplate.works/";
    license = licenses.gpl3;
  };
}
