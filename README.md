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
nix-shell -p '<nixpkgs>' gtk4 gnome3.gtk gobjectIntrospection pkgconfig
```

#### Linux - Ubuntu

```sh
sudo apt install libgtk-3-dev # 18.04 (bionic) or later
sudo apt install libgtk-4-dev # 21.04 (hirsuit) or later
```

#### macOS

```sh
brew install gtk4 gtk+3 gobject-introspection pkg-config
```

### Installing Go

The minimum Go version required to run `gotk4` is 1.17. If your distribution or
operating system does not have 1.17, follow the steps under the "Other
OS/distros" section.

#### Linux/macOS - Nix

```sh
nix-shell # will grab both Go and all GTK dependencies
```

#### Linux - Ubuntu

Snippet taken from [the Go wiki](https://github.com/golang/go/wiki/Ubuntu):

```sh
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install golang-go # or golang-1.17
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
