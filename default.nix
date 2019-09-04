# Copyright (C) 2016-2019  Vincent Ambo <mail@tazj.in>
#
# This file is part of Kontemplate.
#
# Kontemplate is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This file is the Nix derivation used to install Kontemplate on
# Nix-based systems.

{ pkgs ? import <nixpkgs> {} }:

with pkgs; buildGoPackage rec {
  name = "kontemplate-${version}";
  version = "master";
  src = ./.;
  goPackagePath = "github.com/tazjin/kontemplate";
  goDeps = ./deps.nix;
  buildInputs = [ parallel ];

  # Enable checks and configure check-phase to include vet:
  doCheck = true;
  preCheck = ''
    for pkg in $(getGoDirs ""); do
      buildGoDir vet "$pkg"
    done
  '';

  meta = with lib; {
    description = "A resource templating helper for Kubernetes";
    homepage = "http://kontemplate.works/";
    license = licenses.gpl3;
  };
}
