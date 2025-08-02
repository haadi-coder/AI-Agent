package main

import (
	"encoding/json"
	"io/fs"
	"path/filepath"
)

var ListFilesDefinition = ToolDefinition{
	Name:        "list_files",
	Description: "tool to list entries of certain path. If there isnt any path provided, list entries of current directory.",
	InputSchema: GenerateSchema[ListFilesInput](),
	Function:    ListFiles,
}

type ListFilesInput struct {
	Path string `json:"path" jsonschema_description:"relative path of directory to list files. Default value is current directory if any path dont provided"`
}

func ListFiles(input json.RawMessage) (string, error) {
	listFilesInput := ListFilesInput{}
	if err := json.Unmarshal(input, &listFilesInput); err != nil {
		panic(err)
	}

	dir := "."
	if listFilesInput.Path != "" {
		dir = listFilesInput.Path
	}

	var files []string
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if relPath != "." {
			if info.IsDir() {
				files = append(files, relPath+"/")
			} else {
				files = append(files, relPath)
			}
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(files)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
