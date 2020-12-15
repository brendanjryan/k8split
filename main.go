package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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

		dec := yaml.NewDecoder(bytes.NewReader(d))

		names := map[string]int{}
		i := 0

		for {
			data := map[string]interface{}{}
			if err := dec.Decode(&data); err != nil {
				if err.Error() == "EOF" {
					break
				}
				log.Fatalf("error reading yaml document %d: %s", i, err)
			}

			// skip empty documents
			if len(data) < 1 {
				continue
			}

			p, err := yaml.Marshal(data)
			if err != nil {
				log.Fatalf("error creating yaml for document %d: %s", i, err)
			}

			// deduce the name of the
			kind, ok := data["kind"].(string)
			if !ok {
				log.Fatalf("no `Kind` field specified for yaml document %d in this file.", i)
			}

			metadata, ok := data["metadata"].(map[string]interface{})
			if !ok {
				log.Fatalf("no `Metadata` field specified for yaml document %d in this file.", i)
			}

			n, ok := metadata["name"].(string)
			if !ok {
				log.Fatalf("no `Metadata.name` field specified for yaml document %d in this file.", i)
			}

			name := fmt.Sprintf("%s-%s", kind, n)

			c, _ := names[name]
			names[name] = c + 1

			fName := fmt.Sprintf("%s_%d.yaml", strings.ToLower(name), c)
			if c == 0 {
				fName = fmt.Sprintf("%s.yaml", strings.ToLower(name))
			}

			log.Println("Writing file:", fName)

			err = ioutil.WriteFile(fmt.Sprintf("%s/%s", outDir, fName), p, 0644)
			if err != nil {
				log.Fatal("error writing file: ", err)
			}
			i++
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
