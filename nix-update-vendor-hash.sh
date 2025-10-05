#!/usr/bin/env bash
# Helper script to calculate and update vendorHash in flake.nix

set -e

echo "Calculating vendor hash for oneko-go..."
echo

# Calculate the hash
HASH=$(nix-prefetch -f '<nixpkgs>' 'buildGoModule {
  pname = "oneko-go";
  version = "2.0.2";
  src = ./.;
  vendorHash = "";
}' 2>&1 | grep -oP 'sha256-[A-Za-z0-9+/=]+' | head -1 || echo "")

if [ -z "$HASH" ]; then
  echo "Failed to calculate hash. Trying alternative method..."

  # Try building and extract hash from error
  nix build 2>&1 | grep -oP 'sha256-[A-Za-z0-9+/=]+' | head -1 > /tmp/vendor-hash.txt || true

  if [ -s /tmp/vendor-hash.txt ]; then
    HASH=$(cat /tmp/vendor-hash.txt)
  fi
fi

if [ -z "$HASH" ]; then
  echo "Error: Could not calculate vendor hash"
  echo "Please run 'nix build' and check the error message for the hash"
  exit 1
fi

echo "Calculated vendor hash: $HASH"
echo

# Update flake.nix
if [[ "$OSTYPE" == "darwin"* ]]; then
  # macOS
  sed -i '' "s|vendorHash = \"sha256-[A-Za-z0-9+/=]*\";|vendorHash = \"$HASH\";|g" flake.nix
else
  # Linux
  sed -i "s|vendorHash = \"sha256-[A-Za-z0-9+/=]*\";|vendorHash = \"$HASH\";|g" flake.nix
fi

echo "Updated flake.nix with new vendor hash"
echo "New hash: $HASH"
