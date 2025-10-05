# Building Oneko with Nix Flakes

This project includes a Nix flake for reproducible builds on NixOS and other Linux systems with Nix installed.

## Prerequisites

- Nix package manager with flakes enabled
- (Optional) direnv for automatic environment loading

### Enable Flakes

Add to `~/.config/nix/nix.conf` or `/etc/nix/nix.conf`:

```
experimental-features = nix-command flakes
```

## Quick Start

### Build the application

```bash
nix build
```

The binary will be in `./result/bin/oneko-go`

### Run directly

```bash
nix run
```

### Enter development shell

```bash
nix develop
```

This provides:
- Go compiler
- All required X11 and OpenGL libraries
- Development tools (gopls, golangci-lint, etc.)

## Using direnv (Recommended)

If you have direnv installed:

```bash
# Allow the .envrc file (first time only)
direnv allow

# The development environment will automatically load
# when you cd into the project directory
```

## Building Options

### Default build (with fetched dependencies)

```bash
nix build .#default
```

### Vendored build

First, create a vendor directory:

```bash
go mod vendor
```

Then build with vendored dependencies:

```bash
nix build .#oneko-go-vendored
```

## Development Workflow

```bash
# Enter development shell
nix develop

# Build
go build -o neko

# Run
./neko

# Format code
go fmt ./...

# Run linter
golangci-lint run
```

## Installing

### Install to user profile

```bash
nix profile install .#default
```

### Install system-wide (NixOS)

Add to your NixOS configuration:

```nix
{
  inputs.oneko-go.url = "github:glreno/oneko-go";

  # In your system packages:
  environment.systemPackages = [
    inputs.oneko-go.packages.${system}.default
  ];
}
```

Or with flake inputs:

```nix
# flake.nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    oneko-go.url = "github:glreno/oneko-go";
  };

  outputs = { self, nixpkgs, oneko-go }: {
    nixosConfigurations.yourhost = nixpkgs.lib.nixosSystem {
      # ...
      modules = [
        {
          environment.systemPackages = [
            oneko-go.packages.x86_64-linux.default
          ];
        }
      ];
    };
  };
}
```

## Updating Dependencies

If you modify go.mod, you need to update the vendorHash:

```bash
# Try to build (it will fail with the correct hash)
nix build

# Update flake.nix with the hash from the error message
# Replace the vendorHash value with the sha256 from the error
```

Or calculate it manually:

```bash
nix-prefetch -f '<nixpkgs>' 'buildGoModule {
  pname = "oneko-go";
  version = "2.0.2";
  src = ./.;
  vendorHash = "";
}'
```

## Flake Structure

The flake provides:

- **packages.default** - Main oneko-go package
- **packages.oneko-go-vendored** - Version using go mod vendor
- **devShells.default** - Development environment
- **apps.default** - Runnable application
- **formatter** - nixpkgs-fmt for formatting Nix files

## Troubleshooting

### Build fails with "vendor directory not found"

For the vendored build, run:
```bash
go mod vendor
git add vendor
```

### X11 libraries not found

The flake automatically includes all required X11 libraries. If you encounter issues:

```bash
# Clean build cache
nix build --rebuild

# Or try development shell
nix develop
go build
```

### Wrong vendorHash

Update the `vendorHash` in flake.nix with the correct SHA256 from the error message.

Set to `null` to use a vendor directory instead:
```nix
vendorHash = null;
```

## Resources Location

When installed via Nix, resources are located at:
```
$out/share/oneko/resources/
```

The binary is automatically configured to find them.

## Cross-compilation

Build for different systems:

```bash
# For x86_64-linux
nix build .#packages.x86_64-linux.default

# For aarch64-linux
nix build .#packages.aarch64-linux.default
```

## Clean Build

```bash
# Remove build artifacts
nix-collect-garbage

# Remove all build artifacts including cached dependencies
nix-collect-garbage -d
```

## Benefits of Nix Build

✅ **Reproducible** - Same inputs = same outputs
✅ **Isolated** - No system dependencies conflicts
✅ **Declarative** - All dependencies specified in flake.nix
✅ **Cacheable** - Build results can be cached and shared
✅ **Multiple versions** - Can install different versions simultaneously

## See Also

- [Nix Flakes Documentation](https://nixos.wiki/wiki/Flakes)
- [NixOS Manual](https://nixos.org/manual/nixos/stable/)
- [Nix Pills](https://nixos.org/guides/nix-pills/)
