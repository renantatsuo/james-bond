package tools

import (
	"encoding/json"
	"os"
	"strings"
)

var ReadFile = Tool{
	Name:        "readFile",
	Description: "Reads the content of a specified file and returns it as a string.",
	Args:        structToSchema[ReadFileSchema](),
	Fn:          ReadFileFn,
}

type ReadFileSchema struct {
	Path string   `json:"path"`
	_    struct{} `additionalProperties:"false"`
}

func ReadFileFn(input []byte) (string, error) {
	var inputJson ReadFileSchema
	err := json.Unmarshal(input, &inputJson)
	if err != nil {
		return "", err
	}

	path := inputJson.Path
	if !strings.HasPrefix(path, "/") ||
		!strings.HasPrefix(path, "./") {
		path = "./" + path
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
