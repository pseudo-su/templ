package templ

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/pseudo-su/templ/tree"
)

type Templ struct {
	name     *string
	params   *map[string]string
	tree     *tree.Node
	filepath *string
}

func New() Templ {
	return Templ{}
}

func (templ Templ) Params(params map[string]string) Templ {
	templ.params = &params
	return templ
}

func (templ Templ) File(filepath string) Templ {
	templ.filepath = &filepath
	return templ
}

func (templ Templ) Execute() (tree.Node, error) {
	if templ.tree == nil && templ.filepath != nil {
		t, err := ReadFileIntoTree(*templ.filepath)
		if err != nil {
			return nil, err
		}
		if templ.name == nil {
			templ.name = templ.filepath
		}
		templ.tree = &t
	}
	params := *templ.params
	// Setup default values for the templating
	rootTemplate := template.New(*templ.name).Funcs(template.FuncMap{
		"env": func(args ...string) (string, error) {
			input := ""
			for _, arg := range args {
				input += arg
			}
			envVal := os.Getenv(input)
			if envVal == "" {
				return envVal, errors.New("Environment variable " + input + "required")
			}
			return envVal, nil
		},
		"param": func(args ...string) string {
			input := ""
			for _, arg := range args {
				input += arg
			}
			return params[input]
		},
	}).Delims("${", "}")

	// Walk tree and execute templating logic for any string fields
	return tree.ExecuteTreeTemplate(*templ.tree, *rootTemplate)
}

func ReadFileIntoTree(filename string) (tree.Node, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// TODO: this could be map of strings or raw number/string value
	var format tree.FileType
	switch {
	case strings.HasSuffix(filename, ".json"):
		format = tree.JSON
	case strings.HasSuffix(filename, ".yml") || strings.HasSuffix(filename, ".yaml"):
		format = tree.YAML
	case strings.HasSuffix(filename, ".toml"):
		format = tree.TOML
	default:
		return nil, errors.New("unknown format")
	}

	t, err := tree.ReadIntoTree(bytes, format)
	if err != nil {
		return nil, err
	}

	return t, nil
}
