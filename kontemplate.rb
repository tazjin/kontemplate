# Homebrew binary formula for Kontemplate

class Kontemplate < Formula
  desc "Kontemplate - Extremely simple Kubernetes resource templates"
  homepage "https://github.com/tazjin/kontemplate"
  url "https://github.com/tazjin/kontemplate/releases/download/v1.2.0/kontemplate-1.2.0-f8b6ad6-darwin-amd64.tar.gz"
  sha256 "0cda70956b4d4e4944d5760970aaf20d586ef9d62ee90e2d67b70f03ed85e075"
  version "1.2.0-f8b6ad6"

  def install
    bin.install "kontemplate"
  end
end
