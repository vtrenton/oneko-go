# Nix Flake Integration Summary

## Files Added

### Core Nix Files

1. **flake.nix** - Nix flake configuration
   - Defines the package derivation using `buildGoModule`
   - Includes all X11 and OpenGL dependencies
   - Provides development shell with Go tools
   - Supports multiple build variants (default and vendored)

2. **.envrc** - direnv integration
   - Automatically loads Nix development environment
   - Simple one-line configuration: `use flake`

3. **nix-update-vendor-hash.sh** - Vendor hash calculator
   - Helper script to update vendorHash in flake.nix
   - Handles macOS and Linux sed differences

### Documentation

4. **NIX_BUILD.md** - Comprehensive Nix build guide
   - Installation instructions for all methods
   - Development workflow
   - Troubleshooting guide
   - NixOS integration examples

5. **QUICKSTART.md** - Updated with Nix option
   - Added Nix as recommended installation method
   - Links to detailed Nix documentation

6. **.gitignore** - Updated for Nix
   - Ignores `result` symlinks
   - Ignores `.direnv/` cache

## Nix Flake Features

### Packages

- **packages.default** - Standard build with fetched dependencies
  - Uses vendorHash to verify Go module integrity
  - Automatically installs all required system libraries
  - Copies resources to output share directory

- **packages.oneko-go-vendored** - Build with vendored dependencies
  - Uses `go mod vendor` for offline builds
  - Useful for air-gapped environments

### Development Shell

Provides complete development environment with:
- Go compiler
- X11 development libraries (libX11, libXcursor, etc.)
- OpenGL libraries
- Go development tools:
  - gopls (language server)
  - go-tools
  - golangci-lint
  - gotools

### Apps

- **apps.default** - Direct run configuration
  - Enables `nix run` to execute oneko-go
  - No installation required

## Dependencies Managed by Nix

### Build Inputs (Runtime Libraries)
- xorg.libX11
- xorg.libXcursor
- xorg.libXrandr
- xorg.libXinerama
- xorg.libXi
- libGL
- glfw

### Native Build Inputs (Build Tools)
- pkg-config
- go

### Development Tools (Dev Shell Only)
- go-tools
- gopls
- gotools
- golangci-lint

## Usage Examples

### Build and Run (Zero Configuration)

```bash
# From GitHub
nix run github:glreno/oneko-go

# From local directory
nix build
./result/bin/oneko-go
```

### Development

```bash
# Enter dev shell
nix develop

# Build
go build

# Run
./neko
```

### Install

```bash
# User profile
nix profile install .#default

# Run from anywhere
oneko-go
```

## Key Benefits

### 1. Reproducible Builds
- Same inputs always produce same outputs
- All dependencies pinned in flake.lock
- No "works on my machine" issues

### 2. Zero System Pollution
- All dependencies isolated in Nix store
- No system package installation required
- Multiple versions can coexist

### 3. Declarative Dependencies
- All requirements specified in flake.nix
- Transitive dependencies automatically resolved
- Easy to audit and review

### 4. Cross-Platform Ready
- Works on any Linux distribution
- Works on macOS
- NixOS native integration

### 5. Development Environment
- Instant setup with `nix develop`
- All tools included
- Consistent across team members

## VendorHash Management

The vendorHash ensures Go module integrity:

1. **Initial setup**: Set to placeholder or null
2. **First build**: Nix calculates actual hash
3. **Update**: Use error message or `nix-update-vendor-hash.sh`
4. **Lock**: Committed hash ensures reproducibility

Current vendorHash:
```nix
vendorHash = "sha256-lON42gbi/ugLEbUlABQPLhNI4HnHZ7MB17a7SPCzIGo=";
```

## Resource Handling

Resources are installed to a known location:
```
$out/share/oneko/resources/
```

The binary is wrapped to run from this directory:
```nix
postFixup = ''
  wrapProgram $out/bin/oneko-go \
    --chdir $out/share/oneko
'';
```

## Integration Points

### direnv
With `.envrc`, development environment loads automatically:
```bash
cd oneko-go  # Environment activates automatically
go build     # All dependencies available
```

### NixOS System Configuration
```nix
{
  inputs.oneko-go.url = "github:glreno/oneko-go";

  environment.systemPackages = [
    inputs.oneko-go.packages.${system}.default
  ];
}
```

### Home Manager
```nix
{
  inputs.oneko-go.url = "github:glreno/oneko-go";

  home.packages = [
    inputs.oneko-go.packages.${system}.default
  ];
}
```

## Flake Lock

The `flake.lock` file (auto-generated) pins:
- nixpkgs commit
- flake-utils version
- All transitive dependencies

Update with:
```bash
nix flake update
```

## Testing the Flake

```bash
# Check flake validity
nix flake check

# Show flake metadata
nix flake show

# Show flake outputs
nix flake metadata
```

## Build Variants Comparison

| Variant | VendorHash | Offline | Use Case |
|---------|-----------|---------|----------|
| default | Required | No | Standard builds, CI/CD |
| vendored | null | Yes | Air-gapped, vendored repos |

## Future Enhancements

Possible Nix improvements:

1. **Hydra CI Integration** - Automated builds and binary cache
2. **Cross-compilation** - Build for ARM, etc.
3. **NixOS Module** - systemd service configuration
4. **Overlay** - Add to nixpkgs overlay
5. **Flake Templates** - Provide as template for other Go GUI apps

## Performance Notes

- **First build**: Downloads all dependencies (~100MB)
- **Cached build**: Instant (uses Nix store cache)
- **Incremental**: Only rebuilds changed components
- **Binary cache**: Can share builds across machines

## Troubleshooting

### vendorHash mismatch
```bash
./nix-update-vendor-hash.sh
```

### Missing X11 libraries
Not needed! Nix provides all dependencies automatically.

### direnv not loading
```bash
direnv allow
```

## Comparison: Traditional vs Nix Build

| Aspect | Traditional | Nix |
|--------|-------------|-----|
| System deps | Manual install | Automatic |
| Reproducibility | Version drift | Guaranteed |
| Isolation | System-wide | Per-project |
| Setup time | Minutes | Seconds |
| Conflicts | Possible | Impossible |
| Documentation | Required | Self-documenting |

## Conclusion

The Nix flake provides a modern, reproducible, and user-friendly way to build and distribute Oneko-Go. It eliminates the manual dependency management burden while maintaining full transparency and control over the build process.

For NixOS users, this is the recommended installation method. For other Linux distributions with Nix installed, it provides a superior alternative to traditional package management.
