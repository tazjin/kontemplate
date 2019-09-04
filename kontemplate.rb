# Homebrew binary formula for Kontemplate

class Kontemplate < Formula
  desc "Kontemplate - Extremely simple Kubernetes resource templates"
  homepage "https://github.com/tazjin/kontemplate"
  url "https://github.com/tazjin/kontemplate/releases/download/v1.8.0/kontemplate-1.8.0-6c3b299-darwin-amd64.tar.gz"
  sha256 "c541f39ef14f4822ff2f5472a9c1e57f73f0277b6b85bec8e9514640aac311a8"
  version "kontemplate-1.8.0-6c3b299"

  def install
    bin.install "kontemplate"
  end
end
