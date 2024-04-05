# gotk4 Examples

This repository provides examples on using GTK4 and GTK3 using the
[gotk4][gotk4] bindings.

[gotk4]: https://github.com/diamondburned/gotk4

## Examples

GTK4 examples are in [./gtk4](./gtk4). A simple Hello World example is in
[./gtk4/simple](./gtk4/simple).

## Getting Started

Before running any of these examples, you must grab the dependencies in your
operating system or distribution.

### Installing GTK

Below are instructions for getting the needed dependencies for GTK on certain
distros. Most distros and operating systems will be missing from this list, so
PRs are welcomed.

#### Linux/macOS - Nix

```sh
# Nix channels:
nix-shell -p gtk4 gtk3 gobject-introspection pkg-config

# Nix flakes (very slow):
nix shell nixpkgs#gtk4 nixpkgs#gtk3 nixpkgs#gobject-introspection nixpkgs#pkg-config
```

#### Linux - Arch

```sh
sudo pacman -S gtk4 gobject-introspection
```

#### Linux - Ubuntu

```sh
sudo apt install libgtk-3-dev # 18.04 (bionic) or later
sudo apt install libgtk-4-dev # 21.04 (hirsuit) or later
```

#### Linux - Fedora

```sh
sudo dnf install gtk4-devel gobject-introspection-devel
```

#### Linux - openSUSE

```sh
sudo zypper install gtk4-devel gobject-introspection-devel
```

#### macOS

```sh
brew install gtk4 gtk+3 gobject-introspection pkg-config
```

#### Windows - Msys2

```sh
pacman -S mingw-w64-x86_64-toolchain mingw-w64-x86_64-gtk4 mingw-w64-x86_64-gobject-introspection
```

### Installing Go

The minimum Go version required to run `gotk4` is 1.21. If your distribution or
operating system does not have 1.21, follow the steps under the "Other
OS/distros" section.

#### Linux/macOS - Nix

```sh
# Installing just Go:
nix-env -iA nixpkgs.go

# Dropping into a shell with Go:
nix-shell -p go

# Using gotk4-nix's shell, which will also grab GTK dependencies:
nix-shell
```

#### Linux - Ubuntu

Snippet taken from [the Go wiki](https://github.com/golang/go/wiki/Ubuntu):

```sh
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install golang-go # or golang-1.17
```

#### Linux - Fedora

```sh
sudo dnf install golang
```

The [packaged version of Go](https://src.fedoraproject.org/rpms/golang) on
Fedora  is usually a few versions behind upstream on non-rawhide branches,
so you might want to manually install Go using the installation instructions
mentioned below.

#### Linux - openSUSE

```sh
sudo zypper install go
```

#### Other OS/distros

Follow [doc/install](https://go.dev/doc/install) for more information.

### Running Examples

This snippet assumes that you've already cloned this repository down and are
currently inside the repository.

```sh
go run -v ./gtk4/simple
```

**Important:** if you don't run with `-v`, you might start wondering if `go run`
is hung or not. Always keep in mind that building `gotk4` will be very slow at
first, and **the slow building is normal**.
