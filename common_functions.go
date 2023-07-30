package zaje

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// "reflect"

	"github.com/fatih/color"
	highlight "github.com/jessp01/gohighlight"
	"github.com/urfave/cli"
)

var def *highlight.Def

// SynDir path to github.com/jessp01/gohighlight/syntax_files
var SynDir string

// HighlightLexer lexer to use
var HighlightLexer string

// Debug print debug info
var Debug bool

// AddLineNumbers prefix output with line numbers
var AddLineNumbers bool

// RemoveLineNumbers useful when working on an image input
var RemoveLineNumbers bool

var userSynDir = os.Getenv("HOME") + "/.config/zaje/syntax_files"
var globalSynDir = "/etc/zaje/syntax_files"

// PopulateAppMetadata see https://github.com/urfave/cli/blob/v1.22.14/docs/v1/manual.md#customization-1
func PopulateAppMetadata(app *cli.App) {
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
To use {{.Name}} as a cat replacement:
$ {{.Name}} /path/to/file

To replace tail -f:
$ tail -f /path/to/file | {{.Name}} -l server-log -
(- will make {{.Name}} read progressively from STDIN)

AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COPYRIGHT:
   {{.Copyright}}
   {{end}}
`
	app.Usage = "Syntax highlighter to cover all your shell needs"
	app.Version = "0.21.3-2"
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
			Usage:       "Path to lexer files. The `ZAJE_SYNDIR` ENV var is also honoured.\n   If neither is set, " + userSynDir + " will be used.\n",
			EnvVar:      "ZAJE_SYNDIR",
			Destination: &SynDir,
		},
		cli.StringFlag{
			Name: "lexer, l",
			Usage: `config file to use when parsing input. 
   When none is passed, zaje will attempt to autodetect based on the file name or first line of input. 
   You can set the path to lexer files by exporting the ZAJE_SYNDIR ENV var. 
   If not exported, /etc/zaje/highlight will be used.`,
			Destination: &HighlightLexer,
		},
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Run in debug mode.\n",
			Destination: &Debug,
		},
		cli.BoolFlag{
			Name:        "add-line-numbers, ln",
			Usage:       "Add line numbers.\n",
			Destination: &AddLineNumbers,
		},
	}
}

func printDebugInfo() {
	fmt.Println("DEBUG INFO:")
	fmt.Println("Syntax files dir: " + SynDir)
	fmt.Println("Lexer: " + HighlightLexer)
	fmt.Println("DEFINITIONS:")
	fmt.Println(def)
}

func getDefs(filename string, data []byte) []highlight.LineMatch {

	if SynDir == "" {
		if SynDir == "" {
			if stat, err := os.Stat(userSynDir); err == nil && stat.IsDir() {
				SynDir = userSynDir
			} else {
				if stat, err := os.Stat(globalSynDir); err == nil && stat.IsDir() {
					SynDir = globalSynDir
				}
			}
		}
	}

	var defs []*highlight.Def
	lerr, warnings := highlight.ParseSyntaxFiles(SynDir, &defs)
	if lerr != nil {
		log.Fatal(lerr)
	}

	highlight.ResolveIncludes(defs)

	// Always try to auto detect the best lexer
	if def == nil {
		def = highlight.DetectFiletype(defs, filename, bytes.Split(data, []byte("\n"))[0])
	}

	// if a specific lexer was requested by setting the ENV var, try to load it
	if HighlightLexer != "" {
		syntaxFile, lerr := ioutil.ReadFile(SynDir + "/" + HighlightLexer + ".yaml")
		if lerr == nil {
			// Parse it into a `*highlight.Def`
			def, _ = highlight.ParseDef(syntaxFile)
		}
	}

	if Debug {
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
	lastLineNumberLength := len(fmt.Sprint(len(lines)))
	for lineN, l := range lines {
		colN := 0
		if AddLineNumbers {
			fmt.Print(strings.Repeat(" ", lastLineNumberLength-len(fmt.Sprint(lineN+1))))
			color.Set(color.FgYellow)
			fmt.Print(fmt.Sprintf("%d", lineN+1) + " ")
			color.Unset()
		}
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
					color.Set(color.FgHiRed)

				case highlight.Groups["special"]:
					fallthrough
				case highlight.Groups["type.keyword"]:
					fallthrough
				case highlight.Groups["red"]:
					color.Set(color.FgRed)

				case highlight.Groups["constant"]:
					fallthrough
				case highlight.Groups["constant.number"]:
					fallthrough
				case highlight.Groups["constant.bool"]:
					fallthrough
				case highlight.Groups["symbol.brackets"]:
					fallthrough
				case highlight.Groups["identifier.var"]:
					fallthrough
				case highlight.Groups["cyan"]:
					color.Set(color.FgCyan)

				case highlight.Groups["constant.specialChar"]:
					fallthrough
				case highlight.Groups["constant.string.url"]:
					fallthrough
				case highlight.Groups["constant.string"]:
					fallthrough
				case highlight.Groups["magenta"]:
					color.Set(color.FgHiMagenta)

				case highlight.Groups["type"]:
					fallthrough
				case highlight.Groups["symbol.operator"]:
					fallthrough
				case highlight.Groups["symbol.tag.extended"]:
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

// HandleData process input
func HandleData(filename string, data []byte) {
	matches := getDefs(filename, data)
	if matches == nil {
		fmt.Println(string(data))
		return
	}
	colourOutput(matches, data)
}

// DownloadFile helper function to download image file (when input is a remote image)
func DownloadFile(url, fileName string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	// Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
