{
  description = "A basic flake with a shell";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    templ.url = "github:a-h/templ";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    templ,
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {
        system = "x86_64-linux"; # or something else
        config = {allowUnfree = true;};
      };
      templ-pkg = templ.packages.${system}.templ;
    in {
      devShell = pkgs.mkShell {
        nativeBuildInputs = [pkgs.bashInteractive];
        buildInputs = with pkgs; [
          # LSP
          gopls
          htmx-lsp
          tailwindcss-language-server

          # Build
          go
          air
          tailwindcss

          templ-pkg
          vips
          pkg-config
        ];
      };
    });
}
