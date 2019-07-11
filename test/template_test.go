package test

import (
	"testing"
)

func Test_ChainedStringNodes(t *testing.T) {
	template := `
"0": timtam
"1": ${self "0"}
"2": ${self "1"}
"3": ${self "2"}
"4": ${self "3"}
`
	evalAnd(t, template, nil, nil, assertContainsN("timtam", 5))
}
