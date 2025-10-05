# Quick Start Guide - Oneko Go

## Installation Methods

### Option 1: Nix Flakes (Recommended for NixOS/Nix users)

```bash
# Build and run (no dependencies needed!)
nix run github:glreno/oneko-go

# Or build from local directory
nix build
./result/bin/oneko-go

# Or enter development environment
nix develop
```

See [NIX_BUILD.md](NIX_BUILD.md) for detailed Nix instructions.

### Option 2: Traditional Linux Build

```bash
# 1. Install system dependencies
# Ubuntu/Debian:
sudo apt-get update
sudo apt-get install -y libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libgl1-mesa-dev

# Fedora/RHEL:
sudo dnf install -y libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel mesa-libGL-devel

# 2. Build the application
make build

# Or use the automated system dependency installer:
make install-deps-linux
make build
```

## Installation (macOS)

```bash
# No system dependencies needed
make build
```

## Installation (Windows)

```bash
# No system dependencies needed
go build -o neko.exe
```

## Running

```bash
./neko
```

## Usage

- **Click the cat**: Toggle between window mode and free-floating mode
- **Window mode**: Cat runs around inside a bordered window
- **Free mode**: Borderless transparent window follows your mouse across the desktop
- **Move mouse**: Cat will chase it
- **Keep mouse still**: Cat will eventually sit, lick itself, scratch, yawn, and sleep

## Configuration

Create `~/neko.properties` to customize behavior:

```properties
# Window title (in window mode)
windowTitle=My Neko

# How far mouse must move to wake sleeping cat (pixels)
triggerDistance=64

# How fast the cat runs (pixels per frame)
runDistancePerFrame=16

# Mouse offset (where cat aims relative to cursor)
offsetx=-2
offsety=-3

# How close cat needs to get (pixels)
catchDistance=2

# Animation framerates (frames per second)
runFramerate=8
sitFramerate=3
sharpenFramerate=3
scratchFramerate=10

# Sleep delays (milliseconds)
sleepDelay=1200
yawnDelay=1500
surpriseDelay=1000
```

## Troubleshooting

### Build Error: "X11/Xlib.h: No such file or directory"

**Solution**: Install X11 development libraries

```bash
# Ubuntu/Debian
sudo apt-get install libx11-dev

# Fedora/RHEL
sudo dnf install libX11-devel
```

### Build Error: "cannot find package"

**Solution**: Install Go dependencies

```bash
go get github.com/go-gl/glfw/v3.3/glfw
go get github.com/jezek/xgb
```

Or simply run:
```bash
make deps
```

### Cat doesn't appear

1. Check if the process is running: `ps aux | grep neko`
2. Check if resources directory exists with images: `ls resources/images/`
3. Try window mode by clicking where the cat should be

### Cat doesn't follow mouse

- Ensure X11 is running (Linux)
- Check console for errors
- Verify XGB library is properly installed

## File Structure

```
oneko-go/
├── main.go              # Main application
├── neko_controller.go   # Animation logic
├── neko_settings.go     # Configuration management
├── settings.go          # Properties file loader
├── go.mod              # Go module definition
├── Makefile            # Build automation
├── resources/
│   ├── neko.properties # Default configuration
│   └── images/         # Cat animation frames (1.GIF - 32.GIF)
└── README_GO.md        # Full documentation
```

## Development

```bash
# Format code
make fmt

# Run tests
make test

# Check for issues
make vet

# Clean build artifacts
make clean
```

## Next Steps

1. **Read full documentation**: See `README_GO.md`
2. **Understand the port**: See `GO_PORT_SUMMARY.md`
3. **Customize**: Edit `~/neko.properties`
4. **Contribute**: Submit issues/PRs to the repository

Enjoy your desktop cat! 🐱
