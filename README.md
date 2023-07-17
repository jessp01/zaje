```groff
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
