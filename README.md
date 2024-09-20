# moreutils-go
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/gabe565/moreutils)](https://github.com/gabe565/moreutils/releases)
[![Build](https://github.com/gabe565/moreutils/actions/workflows/build.yaml/badge.svg)](https://github.com/gabe565/moreutils/actions/workflows/build.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gabe565/moreutils)](https://goreportcard.com/report/github.com/gabe565/moreutils)

A Go rewrite of [moreutils](http://kitenet.net/~joey/code/moreutils/): A collection of the Unix tools that nobody thought to write long ago when Unix was young.

Some of the original moreutils commands require Perl, so I decided to rewrite them in a language which can compile into a single binary with no dependencies.

## Applets

- **[chronic](docs/chronic.md)**: Runs a command quietly unless it fails
- **[combine](docs/combine.md)**: Combine sets of lines from two files using boolean operations
- **[errno](docs/errno.md)**: Look up errno names and descriptions
- **[ifne](docs/ifne.md)**: Run a command if the standard input is not empty
- **[isutf8](docs/isutf8.md)**: Check whether files are valid UTF-8
- **[mispipe](docs/mispipe.md)**: Pipe two commands, returning the exit status of the first
- **[parallel](docs/moreutils_parallel.md)**: Run multiple jobs at once
- **[pee](docs/pee.md)**: Tee standard input to pipes
- **[sponge](docs/sponge.md)**: Soak up standard input and write to a file
- **[ts](docs/ts.md)**: Timestamp standard input
- **[vidir](docs/vidir.md)**: Edit a directory in your text editor
- **[vipe](docs/vipe.md)**: Insert a text editor into a pipe
- **[zrun](docs/zrun.md)**: Automatically decompress arguments to command

## Installation

### APT (Ubuntu, Debian)

<details>
  <summary>Click to expand</summary>

1. If you don't have it already, install the `ca-certificates` package
   ```shell
   sudo apt install ca-certificates
   ```

2. Add gabe565 apt repository
   ```
   echo 'deb [trusted=yes] https://apt.gabe565.com /' | sudo tee /etc/apt/sources.list.d/gabe565.list
   ```

3. Update apt repositories
   ```shell
   sudo apt update
   ```

4. Install moreutils-go
   ```shell
   sudo apt install moreutils-go
   ```
</details>

### RPM (CentOS, RHEL)

<details>
  <summary>Click to expand</summary>

1. If you don't have it already, install the `ca-certificates` package
   ```shell
   sudo dnf install ca-certificates
   ```

2. Add gabe565 rpm repository to `/etc/yum.repos.d/gabe565.repo`
   ```ini
   [gabe565]
   name=gabe565
   baseurl=https://rpm.gabe565.com
   enabled=1
   gpgcheck=0
   ```

3. Install moreutils-go
   ```shell
   sudo dnf install moreutils-go
   ```
</details>

### AUR (Arch Linux)

<details>
  <summary>Click to expand</summary>

Install [moreutils-go-bin](https://aur.archlinux.org/packages/moreutils-go-bin) with your [AUR helper](https://wiki.archlinux.org/index.php/AUR_helpers) of choice.
</details>

### Homebrew (macOS, Linux)

<details>
  <summary>Click to expand</summary>

Install moreutils-go from [gabe565/homebrew-tap](https://github.com/gabe565/homebrew-tap):
```shell
brew install gabe565/tap/moreutils-go
```
</details>

### GitHub Actions

<details>
  <summary>Click to expand</summary>

This repository can be added to a GitHub Actions workflow to install the applets.

#### Example
```yaml
name: Example

on: push

jobs:
  example:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: gabe565/moreutils@v0
      - run: echo hello world | ts
```

</details>

### Docker

<details>
  <summary>Click to expand</summary>

A Docker image is available at [`ghcr.io/gabe565/moreutils`](https://ghcr.io/gabe565/moreutils)

In this container, all applets are in the root directory.

#### In a Terminal
Some commands can be run directly from a terminal:
```shell
echo hello world | docker run --rm -i ghcr.io/gabe565/moreutils ts
```

#### While Building a Container
If you are building a container and need one of the applets, you can copy them directly to your container during build:
```dockerfile
FROM alpine
COPY --from=ghcr.io/gabe565/moreutils:0 /usr/bin/ts /usr/bin
CMD echo hello world | ts
```

</details>

### Manual Installation

<details>
  <summary>Click to expand</summary>

1. Download and run the [latest release binary](https://github.com/gabe565/moreutils/releases/latest) for your system and architecture.
2. Extract the binary and place it in the desired directory.
3. Run `moreutils install -sr DIRECTORY` to generate symlinks for each command.
</details>

## Differences

My goal is 100% compatability, but there are currently some differences compared to moreutils:

| Applet       | Differences                                                                                                                                                                                                                                                            |
|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **errno**    | Go does not differentiate between errors with the same number like `EAGAIN` and `EWOULDBLOCK`. This causes `errno 35` to always return `EAGAIN`, and `errno EWOULDBLOCK` to return an error. Special cases may be added in the future to handle these more gracefully. |
| **ifdata**   | This command is not yet implemented, but will be included in a future release.                                                                                                                                                                                         |
| **isutf8**   | Unlike moreutils, which prints the expected value range for non-UTF-8 files, the rewrite only logs the offending line, byte, and char.                                                                                                                                 |
| **lckdo**    | Deprecated in moreutils and intentionally not implemented here. It is recommended to use `flock` as a replacement.                                                                                                                                                     |
| **parallel** | The `-l` flag is not yet supported. Also note parallel is not symlinked by default since [GNU Parallel](https://www.gnu.org/software/parallel/) is typically preferred.                                                                                                |
| **pee**      | The flags `--ignore-sigpipe` and `--ignore-write-errors` are not yet supported.                                                                                                                                                                                        |
| **ts**       | The flags `-r`, `-i`, and `-s` are not yet supported. The `-m` flag will trigger a deprecation warning since Go always uses the system's monotonic clock.                                                                                                              |
