# Homebrew binary formula for Kontemplate

class Kontemplate < Formula
  desc "Kontemplate - Extremely simple Kubernetes resource templates"
  homepage "https://github.com/tazjin/kontemplate"
  url "https://github.com/tazjin/kontemplate/releases/download/v1.0.2/kontemplate-1.0.2-f79b261-darwin-amd64.tar.gz"
  sha256 "5a2db5467bc77e4379b5b98f35c9864010f7023ae01a25fb5cda1aede59e021c"
  version "1.0.2-f79b261"

  def install
    bin.install "kontemplate"
  end
end
