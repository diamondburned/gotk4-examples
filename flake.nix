{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    nixpkgs-gotk4.url = "github:NixOS/nixpkgs/nixos-26.05";
    flake-utils.url = "github:numtide/flake-utils";

    gotk4-nix.url = "github:diamondburned/gotk4-nix";
    gotk4-nix.inputs = {
      nixpkgs.follows = "nixpkgs";
      flake-utils.follows = "flake-utils";
    };
  };

  outputs =
    {
      nixpkgs,
      nixpkgs-gotk4,
      gotk4-nix,
      flake-utils,
      ...
    }:

    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = gotk4-nix.lib.mkShell {
          base.pname = "gotk4-examples";
          pkgs = import nixpkgs-gotk4 {
            inherit system;
            overlays = [
              # gotk4-nix.overlays.patchedGo
              gotk4-nix.overlays.patchelf
            ];
          };
          go = pkgs.go;
          inherit (pkgs) gopls gotools;
        };
      }
    );
}
