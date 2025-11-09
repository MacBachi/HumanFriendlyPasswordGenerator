{
  description = "Go development environment for HumanFriendlyPasswordGenerator";

  inputs = {
    # We use unstable to get the latest Go versions
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs = { self, nixpkgs, ... }@inputs:
    let
      # Define supported systems
      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];

      # Helper function to generate outputs for each system
      forAllSystems = f: nixpkgs.lib.genAttrs supportedSystems (system: f system);

    in
    {
      # Define devShells for `nix develop`
      devShells = forAllSystems (system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        {
          default = pkgs.mkShell {
            
            # Packages made available in the shell
            packages = [
              # The Go toolchain (compiler, etc.)
              pkgs.go

              # Go Language Server (for VSCode, Neovim, etc.)
              pkgs.gopls

              # Additional Go tools (like goimports, gorename)
              pkgs.gotools

              # Debugger for Go
              pkgs.delve
            ];

            # This hook runs when entering the shell
            shellHook = ''
              echo ""
              echo "  Go development environment for HumanFriendlyPasswordGenerator is active."
              echo "  Available tools: $(go version), $(gopls version)"
              echo ""
              
              # Set the Go path, useful if Go modules are used outside $HOME
              # (Often not necessary in modern Go projects, but doesn't hurt)
              export GOPATH=$(pwd)/.go
              export GOBIN=$GOPATH/bin
              export PATH=$GOBIN:$PATH
            '';
          };
        });
    };
}

{
  description = "Go development environment for HumanFriendlyPasswordGenerator";

  inputs = {
    # We use unstable to get the latest Go versions
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs = { self, nixpkgs, ... }@inputs:
    let
      # Define supported systems
      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];

      # Helper function to generate outputs for each system
      forAllSystems = f: nixpkgs.lib.genAttrs supportedSystems (system: f system);

    in
    {
      # Define devShells for `nix develop`
      devShells = forAllSystems (system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        {
          default = pkgs.mkShell {
            
            # Packages made available in the shell
            packages = [
              # The Go toolchain (compiler, etc.)
              pkgs.go

              # Go Language Server (for VSCode, Neovim, etc.)
              pkgs.gopls

              # Additional Go tools (like goimports, gorename)
              pkgs.gotools

              # Debugger for Go
              pkgs.delve
            ];

            # This hook runs when entering the shell
            shellHook = ''
              echo ""
              echo "  Go development environment for HumanFriendlyPasswordGenerator is active."
              echo "  Available tools: $(go version), $(gopls version)"
              echo ""
              
              # Set the Go path, useful if Go modules are used outside $HOME
              # (Often not necessary in modern Go projects, but doesn't hurt)
              export GOPATH=$(pwd)/.go
              export GOBIN=$GOPATH/bin
              export PATH=$GOBIN:$PATH
            '';
          };
        });
    };
}
