{
  description = "Hanchond development flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/25.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        go = pkgs.go;
      in
      {
        devShell = pkgs.mkShell {
          buildInputs = [
            pkgs.go
            pkgs.golangci-lint

            # SQL compiler
            pkgs.sqlc

            # Docs dependencies
            pkgs.bun
          ];

          shellHook = ''
            export GOROOT="${go}/share/go"
            # Not setting GOPATH explicitly will make Go use its default location
            export PATH="$HOME/go/bin:$PATH"
          '';
        };
      });
}