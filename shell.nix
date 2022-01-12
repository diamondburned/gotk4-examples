{ systemPkgs ? import <nixpkgs> {} }:

let unstable = import (systemPkgs.fetchFromGitHub {
		owner = "NixOS";
		repo  = "nixpkgs";
		rev   = "3fdd780";
		hash  = "sha256:0df9v2snlk9ag7jnmxiv31pzhd0rqx2h3kzpsxpj07xns8k8dghz";
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

	lib = systemPkgs.lib;

	gtkPkgs =
		if ((systemPkgs.gtk4 or null) != null && lib.versionAtLeast systemPkgs.gtk4.version "4.4.0")
		then systemPkgs
		else unstable;

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

		unstable.go

		# Development tools.
		# gopls
		# goimports
	];

	CGO_ENABLED = "1";

	TMP    = "/tmp";
	TMPDIR = "/tmp";
}
