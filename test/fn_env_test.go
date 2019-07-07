package test

import (
	"os"
	"testing"
)

func TestEnvSimple(t *testing.T) {
	os.Setenv("ENV", "envname")
	template := `ref_env: ${ env "ENV" }`
	evalAnd(t, template, nil, nil, assertContains("envname"))
}

func TestEnvComposite(t *testing.T) {
	os.Setenv("ENV", "envname")
	os.Setenv("STAGE", "stagename")
	template := `ref_env_composite: ${ env "ENV" }-${ env "STAGE" }`
	evalAnd(t, template, nil, nil, assertContains("envname-stagename"))
}

func TestEnvSpaced(t *testing.T) {
	os.Setenv("ENV", "envname")
	os.Setenv("STAGE", "stagename")
	template := `ref_env_indirect: ${ env "E" "N" "V" }-${ env "S" "T" "A" "G" "E" }`
	evalAnd(t, template, nil, nil, assertContains("envname-stagename"))
}
