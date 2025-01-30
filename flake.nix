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
      in {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go_1_23
            gopls
          ];
        };

        packages.default = pkgs.buildGo123Module {
          pname = "nix-search-tv";
          version = self.rev or "unknown";
          src = self;

          # vendorHash = nixpkgs.lib.fakeHash;
          vendorHash = "sha256-n7v0PMPzXEPv8dLVEnCmKhaihmsKjOL8J+je8vxTthM=";

          subPackages = ["cmd/nix-search-tv"];

          meta = {
            description = "A tool integrating television and nix-search packages";
            homepage = "https://github.com/3timeslazy/nix-search-tv";
            mainProgram = "nix-search-tv";
          };
        };
      }
    );
}
