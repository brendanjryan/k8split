package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	outDir string
)

func init() {
	cmd.Flags().StringVarP(&outDir, "outdir", "o", ".", "The name of the directory.")
}

var cmd = &cobra.Command{
	Use:   "k8split -o <dir> <file>",
	Short: "Split a composite yaml file into multiple distinct files",
	Long:  "Split a composite yaml file into multiple distinct files",

	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires an input file")
		}

		_, err := os.Stat(args[0])
		if err != nil {
			return fmt.Errorf("unable to open file %s - %s", args[0], err)
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {

		_, err := os.Stat(outDir)
		if err != nil {
			log.Fatal(err)
		}

		d, err := ioutil.ReadFile(args[0])
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("splitting %s...", args[0])

		// get the line break style for the current OS
		linebreak := "\n"
		windowsLineEnding := bytes.Contains(d, []byte("\r\n"))
		if windowsLineEnding && runtime.GOOS == "windows" {
			linebreak = "\r\n"
		}

		parts := bytes.Split(d, []byte(linebreak+"---"+linebreak))

		if bytes.Equal(parts[len(parts)-1], []byte("")) {
			parts = parts[:len(parts)-1]
		}

		names := map[string]int{}

		log.Printf("split file into %d chunks", len(parts))

		for i, p := range parts {
			data := map[string]interface{}{}
			err := yaml.Unmarshal(p, &data)
			if err != nil {
				log.Fatal("error loading yaml: ", err)
			}

			// deduce the name of the
			kind, ok := data["kind"].(string)
			if !ok {
				log.Fatalf("no `Kind` field specified for the %d'th document in this file.", i)
			}

			metadata, ok := data["metadata"].(map[interface{}]interface{})
			if !ok {
				log.Fatalf("no `Metadata` field specified for the %d'th document in this file.", i)
			}

			n, ok := metadata["name"].(string)
			if !ok {
				log.Fatalf("no `Metadata.name` field specified for the %d'th document in this file.", i)
			}

			name := fmt.Sprintf("%s-%s", kind, n)

			c, _ := names[name]
			names[name] = c + 1

			fName := fmt.Sprintf("%s_%d.yaml", strings.ToLower(name), c)
			if c == 0 {
				fName = fmt.Sprintf("%s.yaml", strings.ToLower(name))
			}

			log.Println("Writing file:", fName)

			err = ioutil.WriteFile(fmt.Sprintf("%s/%s", outDir, fName), append(p, []byte("\n")...), 0644)
			if err != nil {
				log.Fatal("error writing file: ", err)
			}
		}
	},
}

// usage:
// k8split -o <dir <file>
func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
