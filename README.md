# gmake-shim: a workaround for `gmake` command missing from various GNU make packages

# EXAMPLE

```console
$ cd test
$ gmake
Works!
```

# DOWNLOAD

https://github.com/mcandre/gmake-shim/releases

# ABOUT

It may surprise developers that make is unreliable, with many projects implicitly relying on GNU make syntax and semantics, whereas many environments default to BSD make instead. This tends to break the build. Developers interested in cross-platform support may wish to rename these `makefile`s to `GNUmakefile`s and replace `make` commands with `gmake` to disambiguate.

While this would improve the build in BSD environments, ironically this breaks the build in many GNU/Linux environments. Evidently, package maintainers are also unreliable, often omitting the gmake command from their GNU make packages. You're safe with BSD's, RHEL's, macOS/Homebrew, and macOS/MacPorts, but unsafe in Debian's, Alpine, OpenWrt, Void, Windows/Chocolatey, macOS/Xcode, and macOS/Fink. So while GNU make is available in all of these systems, it seems that there isn't a standard way to trigger GNU make-based builds, without sacrificing support for many platforms either way. make just isn't that portable, unfortunately. How do we deal with this situation?

* Vendor-lock to GNU/Linux, GNU make, and running `make`. No thanks.
* Wrap `make` calls in a dispatcher that invoke `make`, `bmake`, `gmake`, etc. according to the particular environment, overcomplicating everything. Non!
* Rewrite the build in terms of pure POSIX make syntax and semantics. This means no `ifdef`, so this isn't a realistic option but for the very simplest of builds. Also, there aren't many ways to statically verify that a makefile avoids GNU or BSD features, so this is a difficult posture to maintain.
* Rewrite the build explicitly in terms of BSD make. Alas, fewer platforms have a BSD make package available compared to GNU make, so you would be vendor-locking yourself out of Windows/Chocolatey, for example. In addition, BSD make may be unavailable as the `bmake` command on BSD's, so you would arrive at the same problem as GNU make!
* Rewrite the build in a tool with more standard behavior across different platforms, such as cmake, autotools, mage, Shake, Gradle, Gulp, or even bash or POSIX sh scripts. Long term, this is by far the most reliable solution. However, rewriting a project's build system can involve a lot of time and effort.

For projects that depend on GNU make, the most immediately practical option for improving cross-platform support is to run build commands as `gmake` and shim environments that omit this command.

In a pinch, soft links, wrapper scripts, and shell aliases can serve, though these are limited to UNIX environments and the exact details will vary across environments, so we leave these to the particular system in question and encourage package maintainers to resolve this. For Windows, the situation is trickier, because there isn't even a standard shell, but an explosion of Command Prompt, PowerShell, and cygwin-like environments, each with different syntaxes and semantics, requiring a multitude of shims to guarantee uniform execution for different configurations.

Fortunately, we can cut this Gordian knot by simply providing binary shims, that reliably run regardless of shell particulars. Hence, gmake-shim.

# RUNTIME REQUIREMENTS

(None)

# BUILDTIME REQUIREMENTS

* [Go](https://golang.org/) 1.9+
* [Mage](https://magefile.org/) (e.g., `go get github.com/magefile/mage`)
* [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) (e.g. `go get golang.org/x/tools/cmd/goimports`)
* [golint](https://github.com/golang/lint) (e.g. `go get github.com/golang/lint/golint`)
* [errcheck](https://github.com/kisielk/errcheck) (e.g. `go get github.com/kisielk/errcheck`)
* [nakedret](https://github.com/alexkohler/nakedret) (e.g. `go get github.com/alexkohler/nakedret`)
* [gox](https://github.com/mitchellh/gox) (e.g. `go get github.com/mitchellh/gox`)
* [zipc](https://github.com/mcandre/zipc) (e.g. `go get github.com/mcandre/zipc/...`)

# INSTALL FROM REMOTE GIT REPOSITORY

```console
$ go get github.com/mcandre/gmake-shim/...
```

(Yes, include the ellipsis as well, it's the magic Go syntax for downloading, building, and installing all components of a package, including any libraries and command line tools.)

# INSTALL FROM LOCAL GIT REPOSITORY

```console
$ mkdir -p $GOPATH/src/github.com/mcandre
$ git clone https://github.com/mcandre/gmake-shim.git $GOPATH/src/github.com/mcandre/gmake-shim
$ cd $GOPATH/src/github.com/mcandre/gmake-shim
$ sh -c 'cd cmd/gmake-shim && go install ./...'
```

# UNIT TEST

```console
$ go test
```

# INTEGRATION TEST

```console
$ mage integrationTest
```

# UNIT + INTEGRATION TEST

```console
$ mage test
```

# LINT

```console
$ mage lint
```

# PORT

```console
$ mage ports
```
