package test

import (
	"os"
	"testing"
)

func TestEnvSimple(t *testing.T) {
	os.Setenv("ENV", "envname")
	template := `test_case: ${ env "ENV" }`
	evalAnd(t, template, nil, nil, assertContains("envname"))
}

func TestEnvComposite(t *testing.T) {
	os.Setenv("ENV", "envname")
	os.Setenv("STAGE", "stagename")
	template := `test_case: ${ env "ENV" }-${ env "STAGE" }`
	evalAnd(t, template, nil, nil, assertContains("envname-stagename"))
}

func TestEnvSpaced(t *testing.T) {
	os.Setenv("ENV", "envname")
	os.Setenv("STAGE", "stagename")
	template := `test_case: ${ env "E" "N" "V" }-${ env "S" "T" "A" "G" "E" }`
	evalAnd(t, template, nil, nil, assertContains("envname-stagename"))
}
