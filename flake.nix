{
  description = "nix-search-tv";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};

        cmdPkg = "cmd/nix-search-tv";

        mkScript = name: text: pkgs.writeShellScriptBin name text;
        scripts = [
          (mkScript "build" "go build -o $DEV_DIR/bin $CMD_DIR")
          (mkScript "run" "$DEV_DIR/bin/nix-search-tv $@ --config $DEV_DIR/config.json")
          (mkScript "print-search" "run print")
          (mkScript "preview-search" "run preview $@")

          (mkScript "test-integrations" "build && NIX_SEARCH_TV_BIN=$DEV_DIR/bin/nix-search-tv go test -v ./integrations/...")

          (mkScript "build-n-tv" "build && print-search | tv --preview 'echo {} | awk \"{ print $1$2 }\" | xargs preview-search'")
          (mkScript "build-n-fzf" "build && print-search | fzf --wrap --preview 'preview-search {1}{2}'")
        ];

      in {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go_1_23
            gopls
            brotli
            television
            fzf
            tmux
          ];

          buildInputs = scripts;

          shellHook = ''
            export PROJECT_ROOT=$(git rev-parse --show-toplevel)
            export DEV_DIR="$PROJECT_ROOT/.dev"
            export CMD_DIR="$PROJECT_ROOT/${cmdPkg}"
          '';
        };

        packages.default = pkgs.buildGo123Module {
          pname = "nix-search-tv";
          version = self.rev or "unknown";
          src = self;

          # vendorHash = nixpkgs.lib.fakeHash;
          vendorHash = "sha256-uzNDhkovlXx0tIgSJ3E08d0TNmktSrlOOe8Iwi4ZfmU=";

          subPackages = [cmdPkg];

          meta = {
            description = "A tool integrating television and nix-search packages";
            homepage = "https://github.com/3timeslazy/nix-search-tv";
            mainProgram = "nix-search-tv";
          };
        };
      }
    );
}
