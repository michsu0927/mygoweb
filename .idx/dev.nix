# To learn more about how to use Nix to configure your environment
# see: https://developers.google.com/idx/guides/customize-idx-env
{ pkgs, ... }: {
  # Which nixpkgs channel to use.
  channel = "stable-23.11"; # or "unstable"
  # Use https://search.nixos.org/packages to find packages
  packages = [
    pkgs.go
    pkgs.nodejs_20
    pkgs.nodePackages.nodemon
    pkgs.go-licenses
    pkgs.python311
    pkgs.gcc
  ];
  # Sets environment variables in the workspace
  env = { 
    PORT = "3000"; 
    CGO_ENABLED = "1"; };

  idx = {
    
    # Search for the extensions you want on https://open-vsx.org/ and use "publisher.id"
    extensions = [
      "golang.go"
    ];
    workspace = {
      onCreate = {
        # Open editors for the following files by default, if they exist:
        default.openFiles = ["main.go"];
      };
    };
    # Enable previews and customize configuration
    previews = {
      enable = true;
      previews = {
        web = {
          #SET PREVIEW PORT
          env = { PORT = "$PORT"; };
          command = [
            "nodemon"
            "--signal" "SIGHUP"
            "-w" "."
            "-e" "go,html"
            "-x" "go run main.go -addr localhost:$PORT" 
          ];
          
          manager = "web";
        };
      };
    };
  };
  
}