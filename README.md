# ZAJE

[![Go Report Card](https://goreportcard.com/badge/github.com/jessp01/zaje)](https://goreportcard.com/report/github.com/jessp01/zaje)
[![GoDoc](https://godoc.org/github.com/jessp01/zaje?status.svg)](http://godoc.org/github.com/jessp01/zaje)
[![AGPLv3](https://img.shields.io/badge/AGPLv3-blue.svg)](https://github.com/jessp01/zaje/blob/master/LICENSE)

`zaje` is a syntax highlighter that aims to cover all your shell colouring needs. It can act as an ad-hoc replacement for `cat` and, with a spot of one-line shell functions `tail` and other friends.

## Motivation

Highlighting output in the shell is hardly a novel idea and its effectiveness is generally agreed to be high:)
There are other tools that provide similar functionality, for instance `supercat` and `grc`. However, with this
project, I was looking to create a tool that can effectively replace `cat`, `tail` and other traditional utils with zero
to very little effort.

### Features

- Supports over a 100 lexers for programming languages, configuration and log formats and UNIX commands (this is done using the
  [highlight Go package](https://github.com/jessp01/gohighlight)
- Can accept input as an argument as well as from an `STDIN` stream
- Can detect the lexer to use based on:
    * The file name (when acting in `cat` mode)
    * The first line of text (so it will usually work nicely when piping as well)
- Supports explicit specification of the lexer to use via a command-line arg and an `ENV` var
- Easily to deploy: since it's a Go CLI app, it's one, statically linked executable with no dynamic deps
- Easily extendable: see [Revising and adding new lexers](https://github.com/jessp01/gohighlight#revising-and-adding-new-lexers) for details

### Installation

Because `zaje` depends on lexers from the `gohighlight` package and also provides some [helper shell
functions](./utils/functions.rc), I've created [install\_zaje.sh](./install_zaje.sh) to handle its deployment.

This is a shell script and does not require Go to be installed. Simply download and invoke with no arguments:

```sh
$ curl https://raw.githubusercontent.com/jessp01/zaje/master/install_zaje.sh > install_zaje.sh
$ ./install_zaje.sh
```

If you run `install_zaje.sh` as a super user, you only need to start a new shell to get all the functionality.
Otherwise, you'll need to source the functions file (see the script's output for instructions).

Being a Golang application, you can also build it yourself with `go` get or fetch a [specific version](https://github.com/jessp01/zaje/releases).
Fetching from the master branch using `go`:

```sh
$ go install github.com/jessp01/zaje
```

If you take this route, you'll need to copy the `highlight/syntax_files` and `utils/functions.rc` manually.

### Installing `super-zaje`

`super-zaje` does everything `zaje` does but provides the additional functionality of extracting text from an image. 
It's a separate binary because it depends on the [gosseract](https://github.com/otiai10/gosseract) which in turn
depends on `libtesseract` and requires its SOs to be available on the machine.

First, install `zaje` using `install_zaje.sh`, and then...

#### Installing on Debian/Ubuntu
```sh
# install deps:
$ sudo apt-get install -y libtesseract-dev libleptonica-dev tesseract-ocr-eng golang-go
```

Most popular Linux distros include the `libtesseract` package but it may be named differently. If the official repos of
your distro of choice do not have it, you can always compile it from source.

#### Installing on Darwin (what people mistakingly refer to as MacOS)
```sh
$ brew install tesseract
```

Aftering installing `tesseract`, simply invoke the below to install `super-zaje`:

```sh
# install super-zaje
$ go install github.com/jessp01/zaje/super-zaje@v0.21.2-3
```

You can then use it thusly:
```sh
$ ~/go/bin/super-zaje -l sh </path/to/img/of/http/url>
```

For example, try:
```sh
$ ~/go/bin/super-zaje "https://github.com/jessp01/zaje/blob/master/testimg/go1.png?raw=true"
```

**NOTE**: `zaje` is capable of detecting the lexer to use based on the first line of text but with images, you'll often
need to help it and specify a designated lexer by passing `-l $NAME` (e.g: `zaje -l sh`, `zaje -l server-log`, etc).


### ASCIInema screencast (Not a video!)

You can copy all text (commands, outputs, etc) straight off the player:)

[![asciicast](https://asciinema.org/a/599719.svg)](https://asciinema.org/a/599719)

[![asciicast](https://asciinema.org/a/597732.svg)](https://asciinema.org/a/597732)

[![asciicast](https://asciinema.org/a/ltEfcN9sILkUFHruwQLn6rDXm.svg)](https://asciinema.org/a/ltEfcN9sILkUFHruwQLn6rDXm)

### Adding and revising lexers

See [Revising and adding new lexers](https://github.com/jessp01/gohighlight#revising-and-adding-new-lexers).

#### Supported specifiers

```yml
statement: will colour the char group green
identifier: will colour the char group blue
special: will colour the char group red
constant.string | constant | constant.number: will colour the char group cyan
constant.specialChar: will colour the char group magenta
type: will colour the char group yellow
comment: high.green will colour the char group bright green
preproc: will colour the char group bright red

```
Specifying the colour names in the YML is also supported, see [df.yaml](https://github.com/jessp01/gohighlight/blob/master/syntax_files/df.yaml) for an exmaple.

If your new lexer doesn't seem to work, run `zaje` with `-d` or `--debug` to get more info.

### Synopsis

```yml
NAME:
   zaje - Syntax highlighter to cover all your shell needs

USAGE:
   zaje [global options] command [command options] [input-file || - ]
   
COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --syn-dir ZAJE_SYNDIR, -s ZAJE_SYNDIR  Path to lexer files. The ZAJE_SYNDIR ENV var is also honoured.
   If neither is set, ~/.config/zaje/syntax_files will be used. [$ZAJE_SYNDIR]

   --lexer value, -l value  config file to use when parsing input. 
   When none is passed, zaje will attempt to autodetect based on the file name or first line of input. 
   You can set the path to lexer files by exporting the ZAJE_SYNDIR ENV var. 
   If not exported, /etc/zaje/highlight will be used.

   --debug, -d  Run in debug mode.

   --help, -h  show help

   --print-version, -V  print only the version

   
EXAMPLES:
To use zaje as a cat replacement:
$ zaje /path/to/file

To replace tail -f:
$ tail -f /path/to/file | zaje -l server-log -
(- will make zaje read progressively from STDIN)

AUTHOR:
   Jesse Portnoy <jesse@packman.io>
   
COPYRIGHT:
   (c) packman.io

```
