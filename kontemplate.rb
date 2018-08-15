# Homebrew binary formula for Kontemplate

class Kontemplate < Formula
  desc "Kontemplate - Extremely simple Kubernetes resource templates"
  homepage "https://github.com/tazjin/kontemplate"
  url "https://github.com/tazjin/kontemplate/releases/download/v1.7.0/kontemplate-1.7.0-511ae92-darwin-amd64.tar.gz"
  sha256 "44910488c0e0480e306cc7b1de564a4c0d39013130f2f6e89bcd9a7401ef6a9a"
  version "kontemplate-1.7.0-511ae92"

  def install
    bin.install "kontemplate"
  end
end
