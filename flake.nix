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
        # vendorHash = null;
        # vendorHash = pkgs.lib.fakeHash;
        vendorHash = "sha256-AXUk6bTJNolTL6QtBVKU9dTndjvhCtyzF9leLfHbuEk=";

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
          # golangci-lint    # linter (`golangci-lint run`), formatter
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
