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
					patches = [
						# cmd/go/internal/work: concurrent ccompile routines
						(builtins.fetchurl "https://github.com/diamondburned/go/commit/4e07fa9fe4e905d89c725baed404ae43e03eb08e.patch")
						# cmd/cgo: concurrent file generation
						(builtins.fetchurl "https://github.com/diamondburned/go/commit/432db23601eeb941cf2ae3a539a62e6f7c11ed06.patch")
					];
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

in gtkPkgs.mkShell {
	buildInputs = with gtkPkgs; [
		glib
		graphene
		gdk-pixbuf
		gnome3.gtk
		gtk4
		vulkan-headers
	];

	nativeBuildInputs = with gtkPkgs; [
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
