{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    flake-compat.url = "https://flakehub.com/f/edolstra/flake-compat/1.tar.gz";
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
      in {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go_1_23
            gopls
            gotools
          ];
        };

        packages.default = pkgs.buildGoModule {
          pname = "nix-search-tv";
          version = self.rev or "unknown";
          src = self;

          vendorHash = "sha256-bModWDH5Htl5rZthtk/UTw/PXT+LrgyBjsvE6hgIePY=";

          meta = {
            description = "A tool integration television and nix-search packages";
            homepage = "https://github.com/3timeslazy/nix-search-tv";
            mainProgram = "nix-search-tv";
          };
        };

        apps = {
          default = "${self.packages.${system}.default}/bin/nix-search-tv";
        };
      }
    );
}
