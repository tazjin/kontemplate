# Homebrew binary formula for Kontemplate

class Kontemplate < Formula
  desc "Kontemplate - Extremely simple Kubernetes resource templates"
  homepage "https://github.com/tazjin/kontemplate"
  url "https://github.com/tazjin/kontemplate/releases/download/v1.1.0/kontemplate-1.1.0-f7ce04e-darwin-amd64.tar.gz"
  sha256 "1dccc80804589f6bd233dd79f52527f2ba8c4fa8d38857c80b2f47b504fe6c04"
  version "1.1.0-f7ce04e"

  def install
    bin.install "kontemplate"
  end
end
