package zaje

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	highlight "github.com/jessp01/gohighlight"
)

func parseDefs(t *testing.T, filename string, data []byte, highlightLexer string) *highlight.Def {
	var currDef *highlight.Def
	synDir := "./highlight/syntax_files"

	var defs []*highlight.Def
	warnings, lerr := highlight.ParseSyntaxFiles(synDir, &defs)
	if lerr != nil {
		t.Errorf("Couldn't get defs from '%s', error: %v\n", synDir, lerr)
	}
	if len(warnings) > 0 {
		t.Errorf("Parsing ended with warnings: '%s'\n", strings.Join(warnings, ";"))
	}

	highlight.ResolveIncludes(defs)

	// Always try to auto detect the best lexer:was
	if currDef == nil {
		currDef = highlight.DetectFiletype(defs, filename, bytes.Split(data, []byte("\n"))[0])
	}

	// if a specific lexer was requested by setting the ENV var, try to load it
	if highlightLexer != "" {
		syntaxFile, lerr := ioutil.ReadFile(synDir + "/" + highlightLexer + ".yaml")
		if lerr == nil {
			// Parse it into a `*highlight.Def`
			def, _ = highlight.ParseDef(syntaxFile)
		}
	}

	if currDef == nil {
		t.Errorf("Found no defs for '%s'\n", filename)
	}

	return currDef
}

func TestInputs(t *testing.T) {

	testInputsDir := "./highlight/test_inputs/*"
	files, err := filepath.Glob(testInputsDir)
	if err != nil {
		t.Errorf("Couldn't open '%s', error: %v\n", testInputsDir, err)
	}
	for _, filename := range files {
		ext := strings.Split(filename, ".")
		data, _ := ioutil.ReadFile(filename)
		NullifyDef()
		fmt.Println("Testing def detection for " + filename)
		def := parseDefs(t, filename, data, "")
		HandleData(filename, data)
		exp := ext[len(ext)-1]
		if def.FileType != exp {
			t.Errorf("\nInput   [%#v]\nExpected[%#v]\nGot     [%#v]\n",
				// string(data), exp, def.FileType)
				filename, exp, def.FileType)
		}
	}
}
