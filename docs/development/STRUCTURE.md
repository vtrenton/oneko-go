# Project Structure

The project has been restructured to follow Go best practices:

```
oneko-go/
├── cmd/
│   └── oneko/           # Main application entry point
│       └── main.go      # Starts the UI
│
├── internal/            # Private application code
│   ├── config/          # Configuration management
│   │   ├── settings.go         # Generic properties loader
│   │   └── neko_settings.go    # Neko-specific settings
│   │
│   ├── controller/      # Game logic (TODO: move controller here)
│   │
│   └── ui/              # User interface
│       └── app.go       # UI implementation
│
├── assets/              # Static resources
│   └── resources/       # Images and config
│       ├── neko.properties
│       └── images/      # 32 GIF animation frames
│
├── flake.nix            # Nix build configuration
├── go.mod               # Go module definition
├── Makefile             # Build automation
└── *.md                 # Documentation

```

## Old Files Removed

- `src/main/java/` - Java source files
- `pom.xml` - Maven build file
- Root-level `.go` files (moved to appropriate packages)

## New Structure Benefits

1. **Cleaner root directory** - Only configuration and documentation
2. **Standard Go layout** - Follows `cmd/` and `internal/` convention
3. **Better organization** - Code grouped by function
4. **Easier to navigate** - Clear separation of concerns

## Building

```bash
# From project root
go build ./cmd/oneko

# Or with make
make build

# Or with nix
nix build
```

## Next Steps

The UI layer (`internal/ui/`) is currently a stub and will be replaced with a proper GUI implementation using Fyne or similar.
