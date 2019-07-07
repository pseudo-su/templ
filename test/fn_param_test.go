package test

import (
	"testing"
)

func TestParamBasic(t *testing.T) {
	template := `
test_case: Hello ${ param "your_name" }. (${ param "stage" })
`
	params := map[string]string{
		"your_name": "Bob",
		"stage":     "prod",
	}
	evalAnd(t, template, &params, nil, assertContains("Hello Bob. (prod)"))
}
