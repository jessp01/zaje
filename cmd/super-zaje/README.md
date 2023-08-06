
# super-zaje

[![CI][badge-build]][build]
[![GoDoc][go-docs-badge]][go-docs]
[![GoReportCard][go-report-card-badge]][go-report-card]
[![License][badge-license]][license]

### Installing `super-zaje`

`super-zaje` does everything `zaje` does but provides the additional functionality of extracting text from an image. 
It's a separate binary because it depends on the [gosseract](https://github.com/otiai10/gosseract) which in turn
depends on `libtesseract` and requires its SOs to be available on the machine.

First, install `zaje` using [install_zaje.sh](https://github.com/jessp01/zaje/blob/master/install_zaje.sh), and then...

#### Installing on Debian/Ubuntu
```sh
# install deps:
$ sudo apt-get install -y libtesseract-dev libleptonica-dev tesseract-ocr-eng golang-go
```

Most popular Linux distros include the `libtesseract` package but it may be named differently. If the official repos of
your distro of choice do not have it, you can always compile it from source.

#### Installing on Darwin (what people mistakenly refer to as MacOS)
```sh
$ brew install tesseract
```

After installing `tesseract`, simply invoke the below to install `super-zaje`:

```sh
# install super-zaje
$ go install github.com/jessp01/zaje/cmd/super-zaje@latest
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


```yml
NAME:
   super-zaje - Syntax highlighter to cover all your shell needs

USAGE:
   super-zaje [global options] command [command options] [input-file || - ]
   
COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --syn-dir ZAJE_SYNDIR, -s ZAJE_SYNDIR  Path to lexer files. The ZAJE_SYNDIR ENV var is also honoured.
   If neither is set, /home/jesse/.config/zaje/syntax_files will be used. [$ZAJE_SYNDIR]

   --lexer value, -l value  config file to use when parsing input. 
   When none is passed, zaje will attempt to autodetect based on the file name or first line of input. 
   You can set the path to lexer files by exporting the ZAJE_SYNDIR ENV var. 
   If not exported, /etc/zaje/highlight will be used.

   --debug, -d  Run in debug mode.

   --add-line-numbers, --ln  Add line numbers.

   --remove-line-numbers, --rln  Remove line numbers.

   --help, -h  show help

   --print-version, -V  print only the version

   
EXAMPLES:
To use super-zaje as a cat replacement:
$ super-zaje /path/to/file

To replace tail -f:
$ tail -f /path/to/file | super-zaje -l server-log -
(- will make super-zaje read progressively from STDIN)

AUTHOR:
   Jesse Portnoy <jesse@packman.io>
   
COPYRIGHT:
   (c) packman.io
```

[license]: ./LICENSE
[badge-license]: https://img.shields.io/github/license/jessp01/zaje.svg
[go-docs-badge]: https://godoc.org/github.com/jessp01/zaje?status.svg
[go-docs]: https://godoc.org/github.com/jessp01/zaje
[go-report-card-badge]: https://goreportcard.com/badge/github.com/jessp01/zaje
[go-report-card]: https://goreportcard.com/report/github.com/jessp01/zaje
[badge-build]: https://github.com/jessp01/zaje/actions/workflows/go.yml/badge.svg
[build]: https://github.com/jessp01/zaje/actions/workflows/go.yml
