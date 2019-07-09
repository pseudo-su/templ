package test

import (
	"testing"
)

func TestSelfBasic(t *testing.T) {
	template := `
values:
  your_name: bob
test_case: Hello ${ self "values.your_name" }
`
	evalAnd(t, template, nil, nil, assertContains("Hello bob"))
}

// TODO
// func TestSelfInvalid(t *testing.T) {
// 	template := `
// values:
//   your_name: bob
// test_case: ${ self "values.not_valid_path" }
// `
// 	evalAnd(t, template, nil, nil, assertContains("Hello bob"))
// }

func TestSelfBasic2(t *testing.T) {
	template := `
values:
  your_name: bob
test_case: ${"Hello"} ${ self "values.your_name" }
`
	evalAnd(t, template, nil, nil, assertContains("Hello bob"))
}

func TestSelfObjectReference(t *testing.T) {
	template := `
values:
  your_name: bob
  object:
    ONE: 1
    TWO: 2
test_case: ${self "values.object"}
`
	// TODO: assert test_case matches { "ONE": 1, "TWO": 2 }
	evalAnd(t, template, nil, nil, assertContainsN("ONE", 2))
}

func TestSelfArrayReference(t *testing.T) {
	template := `
values:
  array:
    - ONE
    - TWO
test_case: ${self "values.array"}
`
	// TODO: assert test_case matches ["ONE", "TWO"]
	evalAnd(t, template, nil, nil, assertContainsN("ONE", 2))
}

func TestSelfSpaced(t *testing.T) {
	template := `
values:
  array:
    - ONE
    - TWO
test_case: ${self "values" "." "array"}
`
	// TODO: assert test_case matches ["ONE", "TWO"]
	evalAnd(t, template, nil, nil, assertContainsN("ONE", 2))
}

func TestSelfComposite(t *testing.T) {
	template := `
stage: prod
values:
  section_prod:
    url: prod.url.com
  section_dev:
    url: dev.url.com
test_case: ${ self "values.section_" (self "stage") }
`
	// TODO: assert test_case matches {"url": "prod.url.com"}
	evalAnd(t, template, nil, nil, assertContainsN("prod.url.com", 2))
}

//
