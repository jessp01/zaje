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

Being a Golang application, you can also build it yourself with `go` get or fetch a [specific version](https://github.com/jessp01/zaje/releases).

Fetching from the master branch using `go`:

```sh
$ go get -u -v github.com/jessp01/zaje
```

### Screencast (Not a video!)

[![asciicast](https://asciinema.org/a/ltEfcN9sILkUFHruwQLn6rDXm.svg)](https://asciinema.org/a/ltEfcN9sILkUFHruwQLn6rDXm)

### Synopsis

```yml
NAME:
   zaje - Syntax highlighter to cover all your shell needs

USAGE:
   zaje [global options] command [command options] [arguments...]

VERSION:
   0.21

DESCRIPTION:
   Highlights text based on regular expressions/strings/characters matching.
   Can operate on files or data sent to STDIN.


AUTHOR:
   Jesse Portnoy <jesse@packman.io>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --syn-dir ZAJE_SYNDIR, -s ZAJE_SYNDIR  Path to lexer files. The ZAJE_SYNDIR ENV var is also honoured.
   If neither is set, /etc/zaje/highlight will be used. [$ZAJE_SYNDIR]
   --lexer value, -l value  config file to use when parsing input. 
   When none is passed, zaje will attempt to autodetect based on the file name or first line of input. 
   You can set the path to lexer files by exporting the ZAJE_SYNDIR ENV var. 
   If not exported, /etc/zaje/highlight will be used.
   --help, -h           show help
   --print-version, -V  print only the version

COPYRIGHT:
   (c) packman.io
```
