# Homebrew binary formula for Kontemplate

class Kontemplate < Formula
  desc "Kontemplate - Extremely simple Kubernetes resource templates"
  homepage "https://github.com/tazjin/kontemplate"
  url "https://github.com/tazjin/kontemplate/releases/download/v1.5.0/kontemplate-1.5.0-c68518d-darwin-amd64.tar.gz"
  sha256 "61cd9ad3e28f52260458b707fd3120afa53e0610213ddd0ab03f489439f6b8a9"
  version "1.5.0-c68518d"

  def install
    bin.install "kontemplate"
  end
end
