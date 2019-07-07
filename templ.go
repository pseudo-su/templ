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
	name       *string
	params     *map[string]string
	tree       *tree.Node
	filepath   *string
	startDelim string
	endDelim   string
}

func New() Templ {
	return Templ{
		startDelim: "${",
		endDelim:   "}",
	}
}

func (templ Templ) Params(params map[string]string) Templ {
	templ.params = &params
	return templ
}

func (templ Templ) Delims(start, end string) Templ {
	templ.startDelim = start
	templ.endDelim = end
	return templ
}

func (templ Templ) Tree(name string, tree tree.Node) Templ {
	templ.name = &name
	templ.tree = &tree
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
			arg := tree.JoinStr(args)
			envVal := os.Getenv(arg)
			if envVal == "" {
				return envVal, errors.New("Environment variable " + arg + "required")
			}
			return envVal, nil
		},
		"param": func(args ...string) string {
			arg := tree.JoinStr(args)
			return params[arg]
		},
	}).Delims(templ.startDelim, templ.endDelim)

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

func DescribeTree(rootNode tree.Node, ctx *tree.WalkingContext) (string, error) {
	return tree.DescribeTree(rootNode, ctx)
}
