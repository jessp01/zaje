# ZAJE

[![CI][badge-build]][build]
[![GoDoc][go-docs-badge]][go-docs]
[![GoReportCard][go-report-card-badge]][go-report-card]
[![License][badge-license]][license]

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
$ go install github.com/jessp01/zaje/cmd@latest
```

If you take this route, you'll need to copy the `highlight/syntax_files` and `utils/functions.rc` manually.

```yml
NAME:
   zaje - Syntax highlighter to cover all your shell needs

USAGE:
   zaje [global options] command [command options] [input-file || - ]
   
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

[license]: ./LICENSE
[badge-license]: https://img.shields.io/github/license/jessp01/zaje.svg
[go-docs-badge]: https://godoc.org/github.com/jessp01/zaje?status.svg
[go-docs]: https://godoc.org/github.com/jessp01/zaje
[go-report-card-badge]: https://goreportcard.com/badge/github.com/jessp01/zaje
[go-report-card]: https://goreportcard.com/report/github.com/jessp01/zaje
[badge-build]: https://github.com/jessp01/zaje/actions/workflows/go.yml/badge.svg
[build]: https://github.com/jessp01/zaje/actions/workflows/go.yml
