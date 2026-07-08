{
  description = "Golang development shell (nix shell)";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    go-overlay.url = "github:purpleclay/go-overlay";
  };

  outputs = inputs:
    inputs.flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import (inputs.nixpkgs) {
          inherit system;
          overlays = [inputs.go-overlay.overlays.default];
        };

        sequincli = pkgs.buildGoModule rec {
          pname = "sequin-cli";
          version = "0.14.6";

          src = pkgs.fetchFromGitHub {
            owner = "sequinstream";
            repo = "sequin";
            rev = "v${version}";
            hash = "sha256-iNd9lXPa8c15Vh5ktj59tuJIp1Jz3nSqSq3ikSytebs=";
          };

          sourceRoot = "${src.name}/cli";
          vendorHash = "sha256-kbVZN4A+eGQzatRlwemI7Tg8g1F7gOBlkKiRmrQesQw=";

          subPackages = ["."];

          ldflags = [
            "-X main.version=${version}"
          ];

          postInstall = ''
            if [ -f "$out/bin/cli" ]; then
              mv $out/bin/cli $out/bin/sequin
            fi
          '';
        };
      in {
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            # golang
            go-bin.versions."1.26.1"
            air
            sqlc
            goose
            gotestfmt
            natscli
            sequincli

            # utils
            jq
            bruno
            lazygit
            resterm
          ];

          shellHook = ''
            export GOBIN=$HOME/go/bin
            export PATH=$GOBIN:$PATH

            echo "Go dev shell activated! Happy coding 🚀."
          '';
        };
      }
    );
}
