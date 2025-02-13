# nix-search-tv

Fuzzy search for NixOS packages.

[![asciicast](https://asciinema.org/a/afNYMXrhoEwwh3wzOK7FbsFtW.svg)](https://asciinema.org/a/afNYMXrhoEwwh3wzOK7FbsFtW)

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

`nix-search-tv` does not do the search by itself, but rather integrates
with other general purpose fuzzy finders, such as [television](https://github.com/alexpasmantier/television) and [fzf](https://github.com/junegunn/fzf)

### Television

Add `nix_channels.toml` file to your television config directory with the content below:

```toml
[[cable_channel]]
name = "nixpkgs"
source_command = "nix-search-tv print"
preview_command = "nix-search-tv preview {}"
```

### fzf

The most straightforward integration might look like:

```sh
alias ns="nix-search-tv print | fzf --preview 'nix-search-tv preview {}'"
```

## Configuration

By default, the configuration file is looked at `$XDG_CONFIG_HOME/nix-search-tv/config.json`

```jsonc
{
  // What indexes to search by default
  //
  // default: nixpkgs
  "indexes": ["nixpkgs", "home-manager"],

  // How often to look for updates and run
  // indexer again
  //
  // default: 1 week (168h)
  "update_interval": "3h2m1s",

  // Where to store the index files
  //
  // default: $XDG_CACHE_HOME/nix-search-tv
  "cache_dir": "path/to/cache/dir",

  // Whether to show the banner when waiting for
  // the indexing
  "enable_waiting_message": true
}
```

## Searchable package registries (so far)

- [Nixpkgs](https://search.nixos.org/packages?channel=unstable) (without NixOS options)
- [Home Manager](https://github.com/nix-community/home-manager)

## Credits

This project was inspired and wouldn't exist without work done by [nix-search](https://github.com/diamondburned/nix-search) contributors.
