package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

var EditFileDefinition = ToolDefinition{
	Name:        "edit_file",
	Description: "Edit a file by providing the file name and the content to be edited",
	InputSchema: GenerateSchema[EditFileInput](),
	Function:    EditFile,
}

type EditFileInput struct {
	Path   string `json:"path" jsonschema:"The path of the file to edit"`
	OldStr string `json:"old_str" jsonschema:"The string to be replaced"`
	NewStr string `json:"new_str" jsonschema:"The new string to replace the old string"`
}

func EditFile(input json.RawMessage) (string, error) {
	editFileInput := EditFileInput{}
	err := json.Unmarshal(input, &editFileInput)
	if err != nil {
		return "", err
	}

	if editFileInput.Path == "" || editFileInput.OldStr == editFileInput.NewStr {
		return "", nil
	}

	fileContent, err := os.ReadFile(editFileInput.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return createFile(editFileInput.Path, editFileInput.NewStr)
		}

		return "", err
	}

	newContent := strings.ReplaceAll(string(fileContent), editFileInput.OldStr, editFileInput.NewStr)

	if newContent == string(fileContent) && editFileInput.OldStr != "" {
		return "", fmt.Errorf("no changes made to the file %s", editFileInput.Path)
	}

	err = os.WriteFile(editFileInput.Path, []byte(newContent), 0644)
	if err != nil {
		return "", err
	}

	return newContent, nil
}

func createFile(filePath, content string) (string, error) {
	dir := path.Dir(filePath)

	if dir != "." {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %v", filePath, err)
	}

	return filePath, nil
}
