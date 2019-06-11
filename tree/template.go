package tree

import (
	"bytes"
	"fmt"
	"text/template"
)

type ResolveContext struct {
	prevPaths [][]string
	path      []string
	curr      Node
	root      Node
	parent    Node
}

// func ResolveInTemplate

func CreateSelfFn(rctx ResolveContext) func(...string) string {
	// rctx := ResolveContext{
	// 	contexts: []WalkingContext{ctx},
	// }
	return func(args ...string) string {
		// when resolving an object replace
		return "<self>"
	}
}

func CreateFileLoadFn(rctx ResolveContext) func(...string) string {
	return func(args ...string) string {
		// when resolving an object replace
		return "<load>"
	}
}

func ResolveFileLoad(filename string) {

}

func ResolveStringNode(rootTemplate template.Template, rctx ResolveContext) error {
	// curr := rctx.curr.(StringNode)
	val := rctx.curr.(StringNode).raw
	tmpl, err := rootTemplate.Funcs(template.FuncMap{
		"self": CreateSelfFn(rctx),
		"load": CreateFileLoadFn(rctx),
	}).Parse(val)
	if err != nil {
		return err
	}
	data := struct {
		stage string
	}{
		stage: "Some string",
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return err
	}

	result := tpl.String()
	fmt.Println(result)
	if result == "<self>" {
		// If this is true then
	}
	if result == "<load>" {

	}
	return nil
}

func ExecuteTreeTemplate(rootNode Node, rootTemplate template.Template) (Node, error) {
	err := WalkTree(rootNode, func(ctx WalkingContext) error {
		switch ctx.curr.(type) {
		case StringNode:
			rctx := ResolveContext{
				prevPaths: [][]string{},
				path:      ctx.path,
				curr:      ctx.curr,
				root:      ctx.root,
				parent:    ctx.parent,
			}
			err := ResolveStringNode(rootTemplate, rctx)
			if err != nil {
				return err
			}
		}
		return nil
	}, nil)

	return rootNode, err
}
