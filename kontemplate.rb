# Homebrew binary formula for Kontemplate

class Kontemplate < Formula
  desc "Kontemplate - Extremely simple Kubernetes resource templates"
  homepage "https://github.com/tazjin/kontemplate"
  url "https://github.com/tazjin/kontemplate/releases/download/v1.3.0/kontemplate-1.3.0-98daa6b-darwin-amd64.tar.gz"
  sha256 "4372e2c0f1249aa43f67e664335c20748d0b394ada1fcceefa3ae786e839ee47"
  version "1.3.0-98daa6b"

  def install
    bin.install "kontemplate"
  end
end
