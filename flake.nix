{
  description = "Oneko - Desktop cat that chases your mouse cursor";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        # X11 and OpenGL dependencies for GLFW
        buildInputs = with pkgs; [
          xorg.libX11
          xorg.libXcursor
          xorg.libXrandr
          xorg.libXinerama
          xorg.libXi
          xorg.libXxf86vm
          libGL
          glfw
        ];

        nativeBuildInputs = with pkgs; [
          pkg-config
          go
          makeWrapper
        ];

      in
      {
        packages = {
          default = pkgs.buildGoModule {
            pname = "oneko-go";
            version = "2.0.2";

            src = ./.;

            # vendorHash needs to be updated after first build
            # Set to null initially, then update with the hash from error message
            vendorHash = "sha256-4cQhS8zirS+XZCMmLQxunOS5YYXcr5shWU6KLW9SPFI=";

            inherit buildInputs nativeBuildInputs;

            # Specify subdirectory for main package
            subPackages = [ "cmd/oneko" ];

            # Copy resources to output
            postInstall = ''
              mkdir -p $out/share/oneko
              cp -r assets/resources $out/share/oneko/
            '';

            # Wrap the binary to find resources
            postFixup = ''
              mv $out/bin/oneko $out/bin/oneko-go
              wrapProgram $out/bin/oneko-go \
                --chdir $out/share/oneko
            '';

            meta = with pkgs.lib; {
              description = "Desktop cat that chases your mouse cursor";
              homepage = "https://github.com/glreno/oneko-go";
              license = licenses.unlicense;
              platforms = platforms.linux;
              mainProgram = "oneko-go";
            };
          };

          # Alternative package without vendorHash (uses vendoring)
          oneko-go-vendored = pkgs.buildGoModule {
            pname = "oneko-go";
            version = "2.0.2";

            src = ./.;

            vendorHash = null;

            inherit buildInputs nativeBuildInputs;

            # Ensure vendor directory exists
            preBuild = ''
              if [ ! -d vendor ]; then
                echo "Error: vendor directory not found. Run 'go mod vendor' first."
                exit 1
              fi
            '';

            postInstall = ''
              mkdir -p $out/share/oneko
              cp -r resources $out/share/oneko/
            '';

            postFixup = ''
              wrapProgram $out/bin/oneko-go \
                --chdir $out/share/oneko
            '';

            meta = with pkgs.lib; {
              description = "Desktop cat that chases your mouse cursor (vendored)";
              homepage = "https://github.com/glreno/oneko-go";
              license = licenses.unlicense;
              platforms = platforms.linux;
              mainProgram = "oneko-go";
            };
          };
        };

        # Development shell
        devShells.default = pkgs.mkShell {
          inherit buildInputs;

          nativeBuildInputs = nativeBuildInputs ++ (with pkgs; [
            go-tools
            gopls
            gotools
            golangci-lint
          ]);

          shellHook = ''
            echo "Oneko Go development environment"
            echo "Go version: $(go version)"
            echo ""
            echo "Available commands:"
            echo "  go build          - Build the application"
            echo "  go run .          - Run the application"
            echo "  go test ./...     - Run tests"
            echo "  golangci-lint run - Run linter"
            echo ""
            echo "Resources directory: ./resources"
          '';
        };

        # Application to run
        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/oneko-go";
        };

        # Formatter
        formatter = pkgs.nixpkgs-fmt;
      }
    );
}
