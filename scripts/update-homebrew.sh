#!/bin/bash
set -e

if [ -z "$1" ]; then
  echo "Usage: $0 <version>"
  echo "Example: $0 0.1.3"
  exit 1
fi

VERSION=$1
REPO="Doarakko/bigburger"
TMP_DIR=$(mktemp -d)

echo "Downloading release assets for v${VERSION}..."
gh release download "v${VERSION}" --repo "$REPO" --pattern "*.tar.gz" --dir "$TMP_DIR"

echo ""
echo "Calculating sha256..."
DARWIN_AMD64=$(shasum -a 256 "$TMP_DIR/bigburger-darwin-amd64.tar.gz" | awk '{print $1}')
DARWIN_ARM64=$(shasum -a 256 "$TMP_DIR/bigburger-darwin-arm64.tar.gz" | awk '{print $1}')
LINUX_AMD64=$(shasum -a 256 "$TMP_DIR/bigburger-linux-amd64.tar.gz" | awk '{print $1}')

echo "darwin-amd64: $DARWIN_AMD64"
echo "darwin-arm64: $DARWIN_ARM64"
echo "linux-amd64:  $LINUX_AMD64"

echo ""
echo "Generating formula..."

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
OUTPUT_FILE="$SCRIPT_DIR/../Formula/bigburger.rb"
mkdir -p "$(dirname "$OUTPUT_FILE")"

cat > "$OUTPUT_FILE" << EOF
class Bigburger < Formula
  desc "Biiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiig burger!"
  homepage "https://github.com/Doarakko/bigburger"
  version "${VERSION}"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/Doarakko/bigburger/releases/download/v#{version}/bigburger-darwin-amd64.tar.gz"
      sha256 "${DARWIN_AMD64}"
    end

    on_arm do
      url "https://github.com/Doarakko/bigburger/releases/download/v#{version}/bigburger-darwin-arm64.tar.gz"
      sha256 "${DARWIN_ARM64}"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/Doarakko/bigburger/releases/download/v#{version}/bigburger-linux-amd64.tar.gz"
      sha256 "${LINUX_AMD64}"
    end
  end

  def install
    bin.install "bigburger"
  end

  test do
    assert_match "bigburger", shell_output("#{bin}/bigburger --help", 2)
  end
end
EOF

echo ""
echo "Generated: $OUTPUT_FILE"
cat "$OUTPUT_FILE"

echo ""
echo "Cleaning up..."
rm -rf "$TMP_DIR"

echo ""
echo "Done! Copy Formula/bigburger.rb to homebrew-tap/Formula/bigburger.rb"
