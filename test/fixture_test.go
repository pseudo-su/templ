package test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/pseudo-su/templ"
	"gotest.tools/assert"
)

func failOnError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestDifferentFiletypeLoad(t *testing.T) {
	os.Setenv("ENV", "envname")
	os.Setenv("STAGE", "stagename")
	jsonTree, err := templ.New().Params(map[string]string{}).File("fixtures/simple.json").Execute()
	failOnError(t, err)
	yamlTree, err := templ.New().Params(map[string]string{}).File("fixtures/simple.yaml").Execute()
	failOnError(t, err)
	tomlTree, err := templ.New().Params(map[string]string{}).File("fixtures/simple.toml").Execute()
	failOnError(t, err)

	one, err := json.Marshal(jsonTree)
	failOnError(t, err)
	two, err := json.Marshal(yamlTree)
	failOnError(t, err)
	three, err := json.Marshal(tomlTree)
	failOnError(t, err)

	t.Log("JSON")
	t.Log(string(one))
	t.Log("YAML")
	t.Log(string(two))
	t.Log("TOML")
	t.Log(string(three))

	assert.DeepEqual(t, one, two)
	assert.DeepEqual(t, two, three)
	// panic("")
}
