{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system: {
      packages.sammler =
        nixpkgs.legacyPackages.${system}.callPackage ./sammler.nix {};

      defaultPackage = self.packages.${system}.sammler;
    });
}














