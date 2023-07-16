package main

import (
    "bytes"
    "fmt"
    "bufio"
    "io/ioutil"
    "os"
    "log"
    "strings"
    //"reflect"

    "github.com/fatih/color"
    "github.com/jessp01/highlight"
)

var def *highlight.Def

func getDefs(filename string, data []byte) []highlight.LineMatch {

    syn_dir := os.Getenv("SYNDIR")
    if syn_dir == "" {
	syn_dir = "/etc/zaje/highlight"
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

    highlight_lexer := os.Getenv("HIGHLIGHT_LEXER");

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
    fi, err := os.Stdin.Stat()
    if err != nil {
	panic(err)
    }

    var filename string

    if fi.Mode() & os.ModeNamedPipe == 0 {
	if len(os.Args) <= 1 {
		log.Fatal("No input file provided. We need a file or STDIN data.")
		return
	}
	filename = os.Args[1]
	data, _ := ioutil.ReadFile(os.Args[1])
	handleData(filename, data)
    } else {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
	    data := scanner.Text()

	    if data == "stop" {
		break
	    }

	    handleData(filename, []byte(data))

	}

	if err := scanner.Err(); err != nil {
	    panic(err)
	}
    }
}
