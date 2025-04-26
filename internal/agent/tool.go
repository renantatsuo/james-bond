package agent

import (
	"github.com/swaggest/jsonschema-go"
)

type Tool struct {
	Name        string
	Description string
	Fn          ToolFn
	Args        map[string]jsonschema.SchemaOrBool
}

type ToolFn func(input []byte) (string, error)

func structToSchema[T any]() map[string]jsonschema.SchemaOrBool {
	reflector := jsonschema.Reflector{}
	v := new(T)

	schema, err := reflector.Reflect(v)
	if err != nil {
		panic(err)
	}

	return schema.Properties
}
