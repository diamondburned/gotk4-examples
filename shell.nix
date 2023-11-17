{ pkgs ? import <nixpkgs> {} }:

let
	gotk4-nix = pkgs.fetchFromGitHub {
		owner  = "diamondburned";
		repo   = "gotk4-nix";
		rev    = "ad91dabf706946c4380d0a105f0937e4e8ffd75f";
		sha256 = "0rkw9k98qy7ifwypkh2fqhdn7y2qphy2f8xjisj0cyp5pjja62im";
	};
in

import "${gotk4-nix}/shell.nix" {
	base = {
		pname = "gotk4-examples";
		version = "dev";
	};
}
