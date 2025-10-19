{
  description = "go_tmux_sessionizer flake";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";

  outputs = { self, nixpkgs, ... }:
    let
      system = builtins.currentSystem;
      pkgs = import nixpkgs { inherit system; };
    in {
      packages.${system}.default = pkgs.stdenv.mkDerivation {
        pname = "go_tmux_sessionizer";
        version = "1.0.0";

        src = ./.;
        buildInputs = [ pkgs.go ]; # build dependencies
        buildPhase = ''
          go build $src
        '';
        installPhase = ''
          mkdir -p $out/bin
          cp go_tmux_sessionizer $out/bin/tms
        '';

        meta = with pkgs.lib; {
          description =
            "tmux sessionizer written in Go. Inspired by the Primeagen's tmux-sessionizer.";
          license = licenses.mit;
          maintainers = with maintainers; [ cfung89 ];
        };
      };
    };
}
