# ExWrap

A general purpose executable wrapper that can turn any application written in any programming language into an exectuable file.

### Why?

Though there are numerous systems available today that attempts to solve the problem of converting applications usually written in scripting languages into executable apps that can easily be distributed (Electron, PyInstaller, etc.), they all usually have one or more of the following problems.

- Targeted at applications written in a particular programming language (E.g. Electron and PyInstaller).
- Not cross platform (You need to be on a similar machine to the target to be able to compile for such platforms).
- Have bloated configurations and/or require too many steps setting up and using.

### What does `ExWrap` do different?

- `ExWrap` takes a different approach by allowing you to convert any application written in any language into an executable without chaning any part of your code or altering how your code works or having to write with a set of fixed APIs.
- By leveraging the power of Go, `ExWrap` is cross-platform. This means you don't have to get a different machine to build for a new set of users. You can build for Windows from your MacBook. All Golang compiler OS/Arch pairs are supported.
- `ExWrap` tries is best to reduce configurations to the bearest minimum while maintaining the most familiar representation of objects you can have today (JSON) for its configuration meaning you usually don't have to learn anything different to start working with `ExWrap`. 
- Despite this simplification, `ExWrap` still takes the do it yourself (DIY) approach to doing this. Meaning you get to decide what and executable is and what it contains.

### Installation

We do not have installation via package managers yet as `ExWrap` is still in `Pre-Alpha` phase. You can download and built it yourself, or download a pre-built release. 

If you are building it yourself, ensure yo run the `scripts/build-pkg.sh` (for Linux and MacOS users) or `scripts/build-pkg.cmd` (for Windows user) to setup the pre-requisites (Don't worry. You don't need an internet connection for this phase.)

> **NOTE:**
> 
> You'll also need to add `ExWrap` to PATH.

### Configuration

The documentation for this is in progress. However, you can check the examples folder for sample configurations.

For now, you can simply consult the `impl/config.go` file to see the available configurations.

The ONLY requirement for `ExWrap` is the existence of the `exwrap.json` file (or whatever name you choose). For simplicity, it might be preferred to always keep this files at the root of the application directory similar to how `composer.json` and `package.json` are being used today. However, the file can be anywhere on your system.

**The only required configuration element is the _`entry_point`_. This is what specifies which command will be run when the application is launched.**

### Running `ExWrap`

Simply run the command `exwrap` from the directory containing the `exwrap.json` file. 

Alternatively, you can run the following to point to another file other than `exwrap.json`:

```sh
exwrap -config my-custom-exwrap-file.json
```

By default, `ExWrap` will create a directory `build` into the directory from which the `exwrap` command was called and the final executable can be found in that directory. However, it is possible to change the target directory by using the `-dir` flag.

```sh
exwrap -dir MyCustomBuildDirectory
```

You can type `exwrap --help` for more.

### Notice

> **NOTICE:**
> 
> - For MacOS, `ExWrap` currently generates an application (`.app`) 
> file while for Windows and Linux, it generates an installer by 
> default. 
> - When the installers are run, they'll automatically 
> install the app into the configured install directory (or 
> `C:\Program Files\<app name>` and `/home/$USERNAME/<app name>` 
> for Windows and Linux respectively if none is configured).

### Contributing

All forms of contribution is welcomed and appreciated. Kindly open an issue feature requests, bugs, pull requests, and/or suggestions. Please star the project to help booster visibility and increase the community.


