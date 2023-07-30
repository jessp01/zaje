package main

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/jessp01/zaje"
	"github.com/otiai10/gosseract/v2"
	"github.com/urfave/cli"
)

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
	app.Version = "0.21.3"
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
			Destination: &zaje.SynDir,
		},
		cli.StringFlag{
			Name: "lexer, l",
			Usage: `config file to use when parsing input. 
   When none is passed, zaje will attempt to autodetect based on the file name or first line of input. 
   You can set the path to lexer files by exporting the ZAJE_SYNDIR ENV var. 
   If not exported, /etc/zaje/highlight will be used.`,
			Destination: &zaje.HighlightLexer,
		},
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Run in debug mode.\n",
			Destination: &zaje.Debug,
		},
		cli.BoolFlag{
			Name:        "add-line-numbers, ln",
			Usage:       "Add line numbers.\n",
			Destination: &zaje.AddLineNumbers,
		},
		cli.BoolFlag{
			Name:        "remove-line-numbers, rln",
			Usage:       "Remove line numbers.\n",
			Destination: &zaje.RemoveLineNumbers,
		},
	}

	app.Action = func(c *cli.Context) error {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		fi, err := os.Stdin.Stat()
		if err != nil {
			panic(err)
		}

		var filename string
		var data []byte
		var resp *http.Response

		if fi.Mode()&os.ModeNamedPipe == 0 {
			if c.NArg() < 1 {
				return errors.New("no input file provided. `zaje` needs a file or data from STDIN")
			}
			filename = c.Args().Get(0)
			httpRegex := regexp.MustCompile("^http(s)?://")
			if httpRegex.Match([]byte(filename)) {
				resp, err = http.Get(filename)
				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()
				data, err = ioutil.ReadAll(resp.Body)
			} else {
				data, _ = ioutil.ReadFile(filename)
			}

			mimeType := http.DetectContentType(data)
			if strings.HasPrefix(mimeType, "image") {
				imgDestination := os.TempDir() + "/" + filepath.Base(filename)
				zaje.DownloadFile(filename, imgDestination)
				client := gosseract.NewClient()
				defer client.Close()

				client.Trim = true
				client.SetImage(imgDestination)
				client.SetLanguage("eng")
				err := client.SetVariable("tessedit_char_whitelist", " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~")

				text, err := client.Text()
				if err != nil {
					panic(err)
				}
				if zaje.RemoveLineNumbers {
					reLineNumber := regexp.MustCompile(`(?m)^\s*\d+\s(.*)`)
					text = reLineNumber.ReplaceAllString(text, `$1`)
				}
				data = []byte(text)
			}

			zaje.HandleData(filename, data)
		} else {
			// if progressive (i.e `tail -f` or ping)
			if c.Args().Get(0) == "-" {
				scanner := bufio.NewScanner(os.Stdin)

				for scanner.Scan() {
					data := scanner.Text()
					zaje.HandleData(filename, []byte(data))
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
				zaje.HandleData(filename, []byte(data))
			}
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
