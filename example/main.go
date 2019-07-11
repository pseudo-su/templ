//go:generate gobin -m -run github.com/mjibson/esc -o base_config.go -pkg main base-datadog.config.yaml
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

			baseConfigBytes, err := FSByte(false, "/base-datadog.config.yaml")
			exitOnError(err)
			baseTree, err := templ.ReadYAMLIntoTree(baseConfigBytes)
			exitOnError(err)
			layerConfigTree, err := templ.ReadFileIntoTree(path)
			exitOnError(err)

			mergedTree, err := templ.MergeTrees(baseTree, layerConfigTree)
			exitOnError(err)

			// TODO: params
			params := map[string]string{
				"service": "my-service-name",
				"stage":   "prod",
			}
			configTree, err := templ.New().Params(params).Tree(path, mergedTree).Execute()
			exitOnError(err)

			configJSON, err := configTree.MarshalJSONIndent("", "  ")
			exitOnError(err)

			outputDir := filepath.Join(cwd, outputDirFlag)
			err = os.MkdirAll(outputDir, os.ModePerm)
			exitOnError(err)

			configOutputFilepath := filepath.Join(outputDir, "output.json")
			fmt.Println("Writing: " + configOutputFilepath)
			err = ioutil.WriteFile(configOutputFilepath, configJSON, 0644)
			exitOnError(err)

			// TODO: don't hard-code these
			monitorNames := []string{"HighLatencyP90", "High4XX", "HighErrors"}

			// Monitors output
			monitorsDir := filepath.Join(cwd, outputDirFlag, "monitors")
			err = os.MkdirAll(monitorsDir, os.ModePerm)
			exitOnError(err)

			for _, name := range monitorNames {
				selector := fmt.Sprintf("datadog.monitor_definitions.%v", name)
				monitorTree, merr := configTree.SelectNode(selector)
				exitOnError(merr)

				monitorJSON, merr := monitorTree.MarshalJSONIndent("", "  ")
				exitOnError(merr)

				fileName := fmt.Sprintf("%v.json", name)
				outputFilepath := filepath.Join(monitorsDir, fileName)
				fmt.Println("Writing: " + outputFilepath)
				err = ioutil.WriteFile(outputFilepath, monitorJSON, 0644)
				exitOnError(err)
			}

			// TODO: don't hard-code these
			timeboardNames := []string{"ServiceHealth"}

			// Monitors output
			timeboardsDir := filepath.Join(cwd, outputDirFlag, "timeboards")
			err = os.MkdirAll(timeboardsDir, os.ModePerm)
			exitOnError(err)

			for _, name := range timeboardNames {
				selector := fmt.Sprintf("datadog.timeboard_definitions.%v", name)
				timeboardTree, err := configTree.SelectNode(selector)
				exitOnError(err)

				timeboardJSON, err := timeboardTree.MarshalJSONIndent("", "  ")
				exitOnError(err)

				fileName := fmt.Sprintf("%v.json", name)
				outputFilepath := filepath.Join(timeboardsDir, fileName)
				fmt.Println("Writing: " + outputFilepath)
				err = ioutil.WriteFile(outputFilepath, timeboardJSON, 0644)

				exitOnError(err)
			}
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
