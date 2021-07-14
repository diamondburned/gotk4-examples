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
				version = "1.17beta1";
				src = builtins.fetchurl {
					url    = "https://golang.org/dl/go1.17rc1.linux-arm64.tar.gz";
					sha256 = "sha256:0kps5kw9yymxawf57ps9xivqrkx2p60bpmkisahr8jl1rqkf963l";
				};
				doCheck = false;
			});
		})
	];
};

in unstable.mkShell {
	buildInputs = with unstable; [
		glib
		graphene
		gdk-pixbuf
		gnome3.gtk
		gtk4
		vulkan-headers
	];

	nativeBuildInputs =
		with unstable;
		let sh = unstable.writeShellScriptBin;
		in [
			# Build/generation dependencies.
			gobjectIntrospection
			pkgconfig
			go

			# Development tools.
			# gopls
			goimports
		];

	CGO_ENABLED = "1";

	TMP    = "/tmp";
	TMPDIR = "/tmp";
}
