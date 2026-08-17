{
  pkgs,
  lib,
  config,
  inputs,
  ...
}:
let
  # https://github.com/golang-migrate/migrate/issues/1279
  go-migrate-pg = pkgs.go-migrate.overrideAttrs (oldAttrs: {
    tags = [ "postgres" ];
  });
in
{
  # https://devenv.sh/packages/
  packages = with pkgs; [
    curl
    git
    gnumake
    go-migrate-pg
  ];

  # https://devenv.sh/languages/
  languages = {
    go = {
      enable = true;
      delve.enable = true;
      lsp.enable = true;
    };

    javascript = {
      enable = true;
      npm.enable = true;
    };
  };

  # https://devenv.sh/scripts/
  scripts.version.exec = ''
    go version
  '';

  # https://devenv.sh/basics/
  enterShell = ''
    version
  '';
}
