package main

import (
    "bytes"
    "fmt"
    "bufio"
    "io/ioutil"
    "os"
    "log"
    "strings"
    "errors"
    "time"
    //"reflect"

    "github.com/fatih/color"
    "github.com/jessp01/highlight"
    "github.com/urfave/cli"
)

var def *highlight.Def
var syn_dir string
var highlight_lexer string

func getDefs(filename string, data []byte) []highlight.LineMatch {

    if syn_dir == "" {
	if syn_dir == "" {
	    // TODO: we probably want to support ~/.config/zaje too
	    syn_dir = "/etc/zaje/highlight"
	}
    }

    var defs []*highlight.Def
    lerr := highlight.ParseSyntaxFiles (syn_dir, &defs)
    if lerr != nil {
	log.Fatal(lerr)
    }
    highlight.ResolveIncludes(defs)

    // Always try to auto detect the best lexer:was
    if def == nil{
	def = highlight.DetectFiletype(defs, filename, bytes.Split(data, []byte("\n"))[0])
    }

    // if a specific lexer was requested by setting the ENV var, try to load it
    if highlight_lexer != "" {
	syntaxFile, lerr := ioutil.ReadFile(syn_dir + "/" + highlight_lexer + ".yaml")
	if lerr == nil {
	    // Parse it into a `*highlight.Def`
	    def, lerr = highlight.ParseDef(syntaxFile)
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
			    // There are more possible groups available than just these ones
			    if group == highlight.Groups["statement"] {
				    color.Set(color.FgGreen)
			    } else if group == highlight.Groups["identifier"] {
				    color.Set(color.FgBlue)
			    } else if group == highlight.Groups["preproc"] {
				    color.Set(color.FgHiRed)
			    } else if group == highlight.Groups["special"] {
				    color.Set(color.FgRed)
			    } else if group == highlight.Groups["constant.string"] {
				    color.Set(color.FgCyan)
			    } else if group == highlight.Groups["constant"] {
				    color.Set(color.FgCyan)
			    } else if group == highlight.Groups["constant.specialChar"] {
				    color.Set(color.FgHiMagenta)
			    } else if group == highlight.Groups["type"] {
				    color.Set(color.FgYellow)
			    } else if group == highlight.Groups["constant.number"] {
				    color.Set(color.FgCyan)
			    } else if group == highlight.Groups["comment"] {
				    color.Set(color.FgHiGreen)
			    } else {
				    color.Unset()
			    }
		    }
		    fmt.Print(string(c))
		    colN++
	    }
	    if group, ok := matches[lineN][colN]; ok {
		    if group == highlight.Groups["default"] || group == highlight.Groups[""] {
			    color.Unset()
		    }
	    }

	    fmt.Print("\n")
    }
}

func handleData(filename string, data []byte){
    matches := getDefs(filename, data)
    if matches == nil {
	    fmt.Println(string(data))
	    return
    }
    colourOutput(matches, data)
}


func main() {

    app := cli.NewApp()
    app.Name = "zaje"
    app.Usage = "Syntax highlighter to cover all your shell needs"
    app.Version = "0.21"
    app.EnableBashCompletion = true
    cli.VersionFlag = cli.BoolFlag{
	Name: "print-version, V",
	Usage: "print only the version",
    }
    app.Compiled = time.Now()
    app.Description = "Highlights text based on regular expressions/strings/characters matching.\n   Can operate on files or data sent to STDIN.\n"
    app.Authors = []cli.Author{
	cli.Author{
	    Name:  "Jesse Portnoy",
	    Email: "jesse@packman.io",
	},
    }
    app.Copyright = "(c) packman.io"
    app.Flags = []cli.Flag {
	cli.StringFlag{
	    Name: "syn-dir, s",
	    Usage: "Path to lexer files. The `ZAJE_SYNDIR` ENV var is also honoured.\n   If neither is set, /etc/zaje/highlight will be used.\n",
	    EnvVar: "ZAJE_SYNDIR",
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
    }

    app.Action = func(c *cli.Context) error {
	fi, err := os.Stdin.Stat()
	if err != nil {
	    panic(err)
	}

	var filename string

	if fi.Mode() & os.ModeNamedPipe == 0 {
	    if c.NArg() < 1 {
		    return errors.New("No input file provided. `zaje` needs a file or data from STDIN.")
	    }
	    filename = c.Args().Get(0)
	    data, _ := ioutil.ReadFile(filename)
	    handleData(filename, data)
	} else {
	    scanner := bufio.NewScanner(os.Stdin)

	    for scanner.Scan() {
		data := scanner.Text()
		handleData(filename, []byte(data))
	    }

	    if err := scanner.Err(); err != nil {
		return err
	    }
	}
	return nil
    }


    err := app.Run(os.Args)
    if err != nil {
	log.Fatal(err)
    }
}
