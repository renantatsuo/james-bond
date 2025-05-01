package tools

import (
	"encoding/json"
	"fmt"
	"os"
)

var WriteFile = Tool{
	Name:        "writeFile",
	Description: "Writes content to a file. Creates a new file if it does not exist.",
	Args:        structToSchema[WriteFileSchema](),
	Fn:          WriteFileFn,
}

type WriteFileSchema struct {
	Path    string   `json:"path"`
	Content string   `json:"content"`
	_       struct{} `additionalProperties:"false"`
}

func WriteFileFn(input []byte) (string, error) {
	var inputJson WriteFileSchema
	err := json.Unmarshal(input, &inputJson)
	if err != nil {
		return "", err
	}

	path := inputJson.Path

	err = os.WriteFile(path, []byte(inputJson.Content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %v", err)
	}

	return "", nil
}
