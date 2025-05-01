package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

var ignore = map[string]struct{}{
	".git":   {},
	"vendor": {},
}

var ListFiles = Tool{
	Name:        "listFiles",
	Description: "Lists the files in the specified path. Defaults to current path. Use it to find the path for a file when not specified.",
	Args:        structToSchema[ListFilesSchema](),
	Fn:          ListFilesFn,
}

type ListFilesSchema struct {
	Path string   `json:"path"`
	_    struct{} `additionalProperties:"false"`
}

func ListFilesFn(input []byte) (string, error) {
	var inputJson ListFilesSchema
	err := json.Unmarshal(input, &inputJson)
	if err != nil {
		return "", err
	}

	basePath := inputJson.Path
	if basePath == "" {
		basePath = "."
	}

	fileList := []string{}
	err = filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		slash := strings.IndexRune(path, '/')
		if slash != -1 {
			prefix := path[:slash]
			_, shouldIgnore := ignore[prefix]

			if shouldIgnore {
				return nil
			}
		}

		relativePath, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}

		if relativePath != "." {
			if info.IsDir() {
				fileList = append(fileList, relativePath+"/")
			} else {
				fileList = append(fileList, relativePath)
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	content, err := json.Marshal(fileList)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
