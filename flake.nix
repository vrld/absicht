{
  description = "Absicht, the mail composer";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs { inherit system; };
      absicht = pkgs.buildGoModule {
        pname = "absicht";
        version = "0.0.1";
        src = ./.;
        vendorHash = "sha256-a5BQFlpg6G0Cjd8kj0KrE7F6r4tE/wAHBVTQXvrWbg0=";

        meta = {
          license = pkgs.lib.licenses.mit;
        };

        buildInputs = [ ];

      };
    in {
      packages = {
        inherit absicht;
        default = absicht;
      };

      devShells.default = pkgs.mkShell {
        hardeningDisable = ["fortify"];
        buildInputs = with pkgs; [
          go
          gopls
          go-tools         # linter (`staticcheck`)
          delve            # debugger
          gdlv             # GUI for delve
        ];

        shellHook = /*bash*/''
          echo "It's dangerous to go alone, take this!"
          echo "  go run main.go"
          echo "  dlv debug --headless main.go --listen=localhost:2345"
          echo "  gdlv connect localhost:2345"
          echo "  staticcheck -h"
        '';
      };
    });
}
