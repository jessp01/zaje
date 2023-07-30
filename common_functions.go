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

	// "reflect"

	"github.com/fatih/color"
	highlight "github.com/jessp01/gohighlight"
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
			if stat, err := os.Stat(os.Getenv("HOME") + "/.config/zaje/syntax_files"); err == nil && stat.IsDir() {
				SynDir = os.Getenv("HOME") + "/.config/zaje/syntax_files"
			} else {
				if stat, err := os.Stat("/etc/zaje/syntax_files"); err == nil && stat.IsDir() {
					SynDir = "/etc/zaje/syntax_files"
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
