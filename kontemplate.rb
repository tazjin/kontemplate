# Homebrew binary formula for Kontemplate

class Kontemplate < Formula
  desc "Kontemplate - Extremely simple Kubernetes resource templates"
  homepage "https://github.com/tazjin/kontemplate"
  url "https://github.com/tazjin/kontemplate/releases/download/v1.4.0/kontemplate-1.4.0-1f373ca-darwin-amd64.tar.gz"
  sha256 "b034cfec6019c973ea7dd30297e1c3434c595e29106e557698ca24f3a4659e5a"
  version "1.4.0-1f373ca"

  def install
    bin.install "kontemplate"
  end
end
