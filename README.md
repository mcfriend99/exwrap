# ExWrap

A general purpose executable wrapper that can turn any application written in any programming language into an exectuable file.

## Why?

Though there are numerous systems available today that attempts to solve the problem of converting applications usually written in scripting languages into executable apps that can easily be distributed (Electron, PyInstaller, etc.), they all usually have one or more of the following problems.

- Targeted at applications written in a particular programming language (E.g. Electron and PyInstaller).
- Not cross platform (You need to be on a similar machine to the target to be able to compile for such platforms).
- Have bloated configurations and/or require too many steps setting up and using.

## What does `ExWrap` do different?

- `ExWrap` takes a different approach by allowing you to convert any application written in any language into an executable without changing any part of your code or altering how your code works or having to write with a set of fixed APIs.
- By leveraging the power of Go, `ExWrap` is cross-platform. This means you don't have to get a different machine to build for a new set of users. You can build for Windows from your MacBook. All (first-class) Golang compiler OS/Arch pairs are supported.
- `ExWrap` tries its best to reduce configurations to the bearest minimum while maintaining the most familiar representation of objects you can have today (JSON) for its configuration meaning you usually don't have to learn anything different to start working with `ExWrap`. 
- Despite this simplification, `ExWrap` still takes the do it yourself (DIY) approach to doing things. Meaning you get to decide what an executable is and what it contains.

## Roadmap

- [x] Generate Installers and Applications for:
  - [x] Windows
  - [x] Linux
  - [x] MacOS (Exe and .APP bundles)
- [x] Support for major compiler architectures:
  - [x] i386
  - [x] AMD64
  - [x] ARM
  - [x] ARM64
- [x] Cross-platform executable generation
- [x] Pre-Install commands
- [x] Post-Install commands

## Installation

We do not have installation via package managers yet as `ExWrap` is still in the `Pre-Alpha` phase. You can download and built it yourself, or download a pre-built release. 

To build it yourself, simply run the `scripts/build.sh` (for Linux and MacOS users) or `scripts/build.cmd` (for Windows user) file. Ensure you have `go` installed and available in user or system PATH.

> **NOTE:**
> 
> You may also need to add `exwrap` to PATH.

## Configuration

The documentation for this is in progress. However, you can check the examples folder for sample configurations.

For now, you can simply consult the [impl/config.go](https://github.com/mcfriend99/exwrap/blob/main/impl/config.go) file to see the available configurations.

The ONLY requirement for `ExWrap` is the existence of the `exwrap.json` file (or whatever name you choose). For simplicity, it might be preferred to always keep this files at the root of the application directory similar to how `composer.json` and `package.json` are being used today. However, the file can be anywhere on your system.

**The only required configuration element is the _`entry_point`_. This is what specifies which command will be run when the application is launched.**

## Running `ExWrap`

Simply run the command `exwrap` from the directory containing the `exwrap.json` file. 

Alternatively, you can run the following to point to another configuration file other than `exwrap.json`.:

```sh
exwrap -config /path/to/my-custom-exwrap-file.json
```

By default, `ExWrap` generates build into the `build` directory at the location from which the `exwrap` command was called. It will create this directory if it doesn't exist. The final executable can be found in that directory. However, it is possible to change the target directory by using the `-dir` flag.

```sh
exwrap -dir MyCustomBuildDirectory
```

You can type `exwrap --help` for more.

## NOTICE

> **Notice for all users**
>
> - When the installers are run, they'll automatically install the app into the configured install directory (or `C:\Program Files\<app name>` and `/home/$USERNAME/<app name>` for Windows and Linux respectively if none is configured).
> - The name of the executable will be the name set in `target_name` or the name of the root directory if `target_name` is not set.

> **Notice For MacOS users:**
> 
> - While `ExWrap` generates an installer Windows and Linux, for MacOS, it does not generate a DMG based installer but rather an executable installer. This is deliberate as there is currently no efficient cross-platform way to programmatically create DMG file on other operating systems without need for users to install extra dependencies which may of their own introduce new bottlenecks.
> - For this reason, `ExWrap` has been enabled with the capability to generate an application (`.app`) file as an opt-in. To enable this, set the darwin > create_app config to `true`.


## Contributing

All forms of contribution is welcomed and appreciated. Kindly open an issue feature requests, bugs, pull requests, and/or suggestions. Please star the project to help booster visibility and increase the community.


