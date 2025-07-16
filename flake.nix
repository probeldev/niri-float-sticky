{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-25.05";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {
        inherit system;
      };
      niri-float-sticky-package = pkgs.callPackage ./package.nix {};
    in {
      packages = rec {
        niri-float-sticky = niri-float-sticky-package;
        default = niri-float-sticky;
      };

      apps = rec {
        niri-float-sticky = flake-utils.lib.mkApp {
          drv = self.packages.${system}.niri-float-sticky;
        };
        default = niri-float-sticky;
      };

      devShells.default = pkgs.mkShell {
        packages = (with pkgs; [
          go
        ]);
      };
    });
}
