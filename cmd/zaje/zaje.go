package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"

	// "reflect"

	"github.com/jessp01/zaje"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	zaje.PopulateAppMetadata(app)

	app.Action = func(c *cli.Context) error {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		fi, err := os.Stdin.Stat()
		if err != nil {
			panic(err)
		}

		var filename string

		if fi.Mode()&os.ModeNamedPipe == 0 {
			if c.NArg() < 1 {
				return errors.New("no input file provided. " + app.Name + " needs a file or data from STDIN")
			}
			filename = c.Args().Get(0)
			data, err := zaje.ReadDataFromFile(filename)
			if err != nil {
				log.Fatal(err)
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
					log.Fatal(err)
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
