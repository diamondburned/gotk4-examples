{ systemPkgs ? import <nixpkgs> {} }:

let unstable = import (systemPkgs.fetchFromGitHub {
		owner  = "NixOS";
		repo   = "nixpkgs";
		rev    = "fbfb79400a08bf754e32b4d4fc3f7d8f8055cf94";
		sha256 = "0pgyx1l1gj33g5i9kwjar7dc3sal2g14mhfljcajj8bqzzrbc3za";
	}) {
		overlays = [
			(self: super: {
				go = super.go.overrideAttrs (old: {
					version = "1.17";
					src = builtins.fetchurl {
						url    = "https://golang.org/dl/go1.17.linux-amd64.tar.gz";
						sha256 = "sha256:0b9p61m7ysiny61k4c0qm3kjsjclsni81b3yrxqkhxmdyp29zy3b";
					};
					doCheck = false;
				});
			})
		];
	};

	choosePkgs = system: if (system) then systemPkgs else unstable;
	lib = systemPkgs.lib;

	goPkgs = choosePkgs
		((systemPkgs.go or null) != null && lib.versionAtLeast systemPkgs.go.version "1.17");

	gtkPkgs = choosePkgs
		((systemPkgs.gtk4 or null) != null && lib.versionAtLeast systemPkgs.gtk4.version "4.4.0");

in systemPkgs.mkShell {
	buildInputs = with gtkPkgs; [
		glib
		graphene
		gdk-pixbuf
		gnome3.gtk
		gtk4
		vulkan-headers
	];

	nativeBuildInputs = with systemPkgs; [
		# Build/generation dependencies.
		gobjectIntrospection
		pkgconfig
		goPkgs.go

		# Development tools.
		# gopls
		# goimports
	];

	CGO_ENABLED = "1";

	TMP    = "/tmp";
	TMPDIR = "/tmp";
}
