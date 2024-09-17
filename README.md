# moreutils
[![Build](https://github.com/gabe565/moreutils/actions/workflows/build.yaml/badge.svg)](https://github.com/gabe565/moreutils/actions/workflows/build.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gabe565/moreutils)](https://goreportcard.com/report/github.com/gabe565/moreutils)

Go rewrite of [moreutils](http://kitenet.net/~joey/code/moreutils/).

## Applets

- **[chronic](docs/chronic.md):** runs a command quietly unless it fails
- **[combine](docs/combine.md):** combine the lines in two files using boolean - operations
- **[ifne](docs/ifne.md):** run a command if the standard input is not empty
- **[mispipe](docs/mispipe.md):** pipe two commands, returning the exit status of - the first
- **[pee](docs/pee.md):** tee standard input to pipes
- **[sponge](docs/sponge.md):** soak up standard input and write to a file
- **[ts](docs/ts.md):** timestamp standard input
- **[vidir](docs/vidir.md):** edit a directory in your text editor
- **[vipe](docs/vipe.md):** insert a text editor into a pipe
- **[zrun](docs/zrun.md):** automatically uncompress arguments to command

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

### Manual Installation

<details>
  <summary>Click to expand</summary>

1. Download and run the [latest release binary](https://github.com/gabe565/moreutils-go/releases/latest) for your system and architecture.
2. Extract the binary and place it in the desired directory.
3. Run `moreutils install -sr DIRECTORY` to generate symlinks for each command.
</details>
