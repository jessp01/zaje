package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
	//"reflect"

	"github.com/fatih/color"
	"github.com/jessp01/gohighlight"
	"github.com/urfave/cli"
)

var def *highlight.Def
var syn_dir string
var highlight_lexer string
var debug bool

func printDebugInfo() {
	fmt.Println("DEBUG INFO:")
	fmt.Println("Syntax files dir: " + syn_dir)
	fmt.Println("Lexer: " + highlight_lexer)
	fmt.Println("DEFINITIONS:")
	fmt.Println(def)
}

func getDefs(filename string, data []byte) []highlight.LineMatch {

	if syn_dir == "" {
		if syn_dir == "" {
			if stat, err := os.Stat(os.Getenv("HOME") + "/.config/zaje/syntax_files"); err == nil && stat.IsDir() {
				syn_dir = os.Getenv("HOME") + "/.config/zaje/syntax_files"
			} else {
				if stat, err := os.Stat("/etc/zaje/syntax_files"); err == nil && stat.IsDir() {
					syn_dir = "/etc/zaje/syntax_files"
				}
			}
		}
	}

	var defs []*highlight.Def
	lerr, warnings := highlight.ParseSyntaxFiles(syn_dir, &defs)
	if lerr != nil {
		log.Fatal(lerr)
	}

	highlight.ResolveIncludes(defs)

	// Always try to auto detect the best lexer
	if def == nil {
		def = highlight.DetectFiletype(defs, filename, bytes.Split(data, []byte("\n"))[0])
	}

	// if a specific lexer was requested by setting the ENV var, try to load it
	if highlight_lexer != "" {
		syntaxFile, lerr := ioutil.ReadFile(syn_dir + "/" + highlight_lexer + ".yaml")
		if lerr == nil {
			// Parse it into a `*highlight.Def`
			def, _ = highlight.ParseDef(syntaxFile)
		}
	}

	if debug {
		printDebugInfo()
		if len(warnings) > 0 {
			fmt.Println(warnings)
		}
	}

	if def == nil {
		return nil
	}

	h := highlight.NewHighlighter(def)

	return h.HighlightString(string(data))
}

func colourOutput(matches []highlight.LineMatch, data []byte) {
	lines := strings.Split(string(data), "\n")
	for lineN, l := range lines {
		colN := 0
		for _, c := range l {
			if group, ok := matches[lineN][colN]; ok {
				switch group {
				case highlight.Groups["default"]:
					fallthrough
				case highlight.Groups[""]:
					color.Unset()
				case highlight.Groups["statement"]:
					fallthrough
				case highlight.Groups["green"]:
					color.Set(color.FgGreen)

				case highlight.Groups["identifier"]:
					fallthrough
				case highlight.Groups["blue"]:
					color.Set(color.FgHiBlue)

				case highlight.Groups["preproc"]:
					//fallthrough
					//case highlight.Groups["high.red"]:
					color.Set(color.FgHiRed)

				case highlight.Groups["special"]:
					fallthrough
				case highlight.Groups["red"]:
					color.Set(color.FgRed)

				case highlight.Groups["constant.string"]:
					fallthrough
				case highlight.Groups["constant"]:
					fallthrough
				case highlight.Groups["constant.number"]:
					fallthrough
				case highlight.Groups["symbol.operator"]:
					fallthrough
				case highlight.Groups["symbol.brackets"]:
					fallthrough
				case highlight.Groups["cyan"]:
					color.Set(color.FgCyan)

				case highlight.Groups["constant.specialChar"]:
					fallthrough
				case highlight.Groups["identifier.var"]:
					fallthrough
				case highlight.Groups["magenta"]:
					color.Set(color.FgHiMagenta)

				case highlight.Groups["type"]:
					fallthrough
				case highlight.Groups["yellow"]:
					color.Set(color.FgYellow)

				case highlight.Groups["comment"]:
					fallthrough
				case highlight.Groups["high.green"]:
					color.Set(color.FgHiGreen)
				default:
					color.Unset()
				}
			}
			fmt.Print(string(c))
			colN++
		}

		color.Unset()
		fmt.Print("\n")
	}
}

func handleData(filename string, data []byte) {
	matches := getDefs(filename, data)
	if matches == nil {
		fmt.Println(string(data))
		return
	}
	colourOutput(matches, data)
}

func main() {

	app := cli.NewApp()
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[input-file || - ]{{end}}
   {{if len .Authors}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}{{ "\n" }}
   {{end}}{{end}}{{if .Copyright }}
EXAMPLES:
To use zaje as a cat replacement:
$ zaje /path/to/file

To replace tail -f:
$ tail -f /path/to/file | zaje -l server-log -
(- will make zaje read progressively from STDIN)

AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COPYRIGHT:
   {{.Copyright}}
   {{end}}
`
	app.Name = "zaje"
	app.Usage = "Syntax highlighter to cover all your shell needs"
	app.Version = "0.21.1-9"
	app.EnableBashCompletion = true
	cli.VersionFlag = cli.BoolFlag{
		Name:  "print-version, V",
		Usage: "print only the version",
	}
	app.Compiled = time.Now()
	app.Description = "Highlights text based on regular expressions/strings/characters matching.\n   Can operate on files or data sent to STDIN.\n"
	app.Authors = []cli.Author{
		{
			Name:  "Jesse Portnoy",
			Email: "jesse@packman.io",
		},
	}
	app.Copyright = "(c) packman.io"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "syn-dir, s",
			Usage:       "Path to lexer files. The `ZAJE_SYNDIR` ENV var is also honoured.\n   If neither is set, ~/.config/zaje/syntax_files will be used.\n",
			EnvVar:      "ZAJE_SYNDIR",
			Destination: &syn_dir,
		},
		cli.StringFlag{
			Name: "lexer, l",
			Usage: `config file to use when parsing input. 
   When none is passed, zaje will attempt to autodetect based on the file name or first line of input. 
   You can set the path to lexer files by exporting the ZAJE_SYNDIR ENV var. 
   If not exported, /etc/zaje/highlight will be used.`,
			Destination: &highlight_lexer,
		},
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Run in debug mode.\n",
			Destination: &debug,
		},
	}

	app.Action = func(c *cli.Context) error {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		fi, err := os.Stdin.Stat()
		if err != nil {
			panic(err)
		}

		var filename string

		if fi.Mode()&os.ModeNamedPipe == 0 {
			if c.NArg() < 1 {
				return errors.New("No input file provided. `zaje` needs a file or data from STDIN.")
			}
			filename = c.Args().Get(0)
			data, _ := ioutil.ReadFile(filename)
			handleData(filename, data)
		} else {
			// if progressive (i.e `tail -f` or ping)
			if c.Args().Get(0) == "-" {
				scanner := bufio.NewScanner(os.Stdin)

				for scanner.Scan() {
					data := scanner.Text()
					handleData(filename, []byte(data))
				}

				if err := scanner.Err(); err != nil {
					return err
				}
				// read everything and process
			} else {
				data, err := io.ReadAll(os.Stdin)
				if err != nil {
					panic(err)
				}
				handleData(filename, []byte(data))
			}
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
