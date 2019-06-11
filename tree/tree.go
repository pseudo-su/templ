package tree

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/BurntSushi/toml"
	"github.com/go-yaml/yaml"
)

type FileType uint

const (
	YAML FileType = iota
	JSON
	TOML
)

func unmarshal(bytes []byte, format FileType) (interface{}, error) {
	var parsedBytes interface{}
	switch format {
	case JSON:
		err := json.Unmarshal(bytes, &parsedBytes)
		if err != nil {
			return nil, err
		}
	case YAML:
		err := yaml.Unmarshal(bytes, &parsedBytes)
		if err != nil {
			return nil, err
		}
	case TOML:
		err := toml.Unmarshal(bytes, &parsedBytes)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("not a valid file suffix")
	}
	return parsedBytes, nil
}

func ReadIntoTree(bytes []byte, format FileType) (Node, error) {
	parsedVal, err := unmarshal(bytes, format)
	if err != nil {
		return nil, err
	}
	v := reflect.ValueOf(parsedVal)
	tree, err := NewNode(v)
	if err != nil {
		return nil, err
	}
	return tree, nil
}
