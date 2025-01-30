# nix-search-tv

An integration between [television](https://github.com/alexpasmantier/television) and [nix-search](https://github.com/diamondburned/nix-search)

[![asciicast](https://asciinema.org/a/AUt4rfSukwSWsrlis7ZNsBP4N.svg)](https://asciinema.org/a/AUt4rfSukwSWsrlis7ZNsBP4N)

## Installation

### Flake

There are many ways how one can install a package from a flake, below is one:

```nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    nix-search-tv.url = "github:3timeslazy/nix-search-tv";
  };

  outputs = {
    nixpkgs,
    nix-search-tv,
    ...
  }: {
    nixosConfigurations.system = nixpkgs.lib.nixosSystem {
      modules = [
        {
          environment.systemPackages = [
            nix-search-tv.packages.x86_64-linux.default
          ];
        }
      ];
    };
  };
}
```

### Go

```sh
git clone https://github.com/3timeslazy/nix-search-tv
cd nix-search-tv
go install ./cmd/nix-search-tv

# `go install github.com/3timeslazy/nix-search-tv/cmd/nix-search-tv@latest` won't work
# because go.mod contains a replace directive
```

## Usage

Once nix-search-tv is installed, first index the Nixpkgs:

```sh
nix-search-tv index

# or

nix-search-tv index --flake nixpkgs
```

Then, add `nix_channels.toml` file to your television config directory with the content below:

```toml
[[cable_channel]]
name = "nixpkgs"
source_command = "nix-search-tv print"
preview_command = "nix-search-tv preview {}"
```

## Credits

This project was inspired and wouldn't exist without work done by [nix-search](https://github.com/diamondburned/nix-search) contributors.
