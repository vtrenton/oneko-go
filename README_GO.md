# Oneko-Go

This is a Go port of the Java Neko desktop pet application.

## Prerequisites

To build this application, you need:

### Linux
```bash
# Debian/Ubuntu
sudo apt-get install libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libgl1-mesa-dev

# Fedora/RHEL
sudo dnf install libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel mesa-libGL-devel
```

### macOS
```bash
# GLFW will use Cocoa framework (no additional dependencies needed)
```

### Windows
```bash
# No additional dependencies needed (uses Win32 API)
```

## Building

```bash
go build -o neko
```

## Running

```bash
./neko
```

## Configuration

Create a `~/neko.properties` file to customize settings. See `resources/neko.properties` for available options.

## Architecture

The Go port maintains the same structure as the Java version:

- **settings.go** - Configuration file loader (replaces Settings.java)
- **neko_settings.go** - Neko-specific settings (replaces NekoSettings.java)
- **neko_controller.go** - Animation and behavior controller (replaces NekoController.java)
- **main.go** - Main application and window management (replaces Neko.java)

## Key Differences from Java Version

1. Uses GLFW instead of Java Swing for cross-platform window management
2. Uses XGB for X11 mouse tracking on Linux
3. Native Go concurrency (goroutines) instead of Java Timers
4. Embedded resources or filesystem-based resource loading

## License

This is public domain software under the UNLICENSE.
http://unlicense.org
