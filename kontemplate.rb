# Homebrew binary formula for Kontemplate

class Kontemplate < Formula
  desc "Kontemplate - Extremely simple Kubernetes resource templates"
  homepage "https://github.com/tazjin/kontemplate"
  url "https://github.com/tazjin/kontemplate/releases/download/v1.6.0/kontemplate-1.6.0-97bef90-darwin-amd64.tar.gz"
  sha256 "d21529153d369d2347477f981a855525695c0bc2912a50f05d3baf15c96f7c16"
  version "1.6.0-97bef90"

  def install
    bin.install "kontemplate"
  end
end
