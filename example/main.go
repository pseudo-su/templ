//go:generate gobin -m -run github.com/mjibson/esc -o base_config.go -pkg main base-config.yaml
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"path/filepath"

	"github.com/pseudo-su/templ"
	"github.com/spf13/cobra"
)

func exitOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	rootCwd, osErr := os.Getwd()
	exitOnError(osErr)
	var filepathFlag string
	var cwdFlag string
	var outputDirFlag string
	var cmdPackage = &cobra.Command{
		Use:   "package [generate files into .stencil folder]",
		Short: "Package config files into .stencil/",
		Long: `package is for generating service config.
Files will be generated into your .stencil/ folder.`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cwd := filepath.Join(rootCwd, cwdFlag)
			path := filepath.Join(cwd, filepathFlag)
			fmt.Println("Loading: " + path)

			baseConfigBytes, err := FSByte(false, "/base-config.yaml")
			exitOnError(err)
			baseTree, err := templ.ReadYAMLIntoTree(baseConfigBytes)
			exitOnError(err)
			configTree, err := templ.ReadFileIntoTree(path)
			exitOnError(err)

			mergedTree, err := templ.MergeTrees(baseTree, configTree)
			exitOnError(err)

			// TODO: params
			params := map[string]string{
				// "service": "service-name",
				// "stage":   "test",
			}
			resultTree, err := templ.New().Params(params).Tree(path, mergedTree).Execute()
			exitOnError(err)

			resultTreeDesc, err := templ.DescribeTree(resultTree)
			exitOnError(err)

			outputFilepath := filepath.Join(cwd, outputDirFlag, "output.txt")

			fmt.Println("Writing: " + outputFilepath)
			err = os.MkdirAll(filepath.Dir(outputFilepath), os.ModePerm)
			exitOnError(err)
			err = ioutil.WriteFile(outputFilepath, []byte(resultTreeDesc), 0644)

			exitOnError(err)
		},
	}

	cmdPackage.Flags().StringVarP(&filepathFlag, "file", "f", "config.yaml", "config file to load")
	cmdPackage.Flags().StringVarP(&outputDirFlag, "outputDir", "o", ".stencil", "package output folder")
	cmdPackage.Flags().StringVarP(&cwdFlag, "cwd", "d", ".", "folder context to execute command in")

	var rootCmd = &cobra.Command{Use: "stencil"}
	rootCmd.AddCommand(cmdPackage)

	cmdErr := rootCmd.Execute()
	exitOnError(cmdErr)
}
