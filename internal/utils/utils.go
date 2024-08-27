package utils

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"masscode-parser/internal/templates"
	"masscode-parser/internal/types"
)

func EnsureCorrectExtension(path string, outputType string) string {
	// Remove any existing extension
	path = strings.TrimSuffix(path, filepath.Ext(path))

	// Append the correct extension based on the output type
	extension := getFileExtension(outputType)
	return path + "." + extension
}

func ParseDatabase(filePath string) (*types.Database, error) {
	var db *types.Database
	data, err := os.ReadFile(filePath)
	if err != nil {
		return db, fmt.Errorf("error reading file: %v", err)
	}

	err = json.Unmarshal(data, &db)
	if err != nil {
		return db, fmt.Errorf("error parsing JSON: %v", err)
	}

	return db, nil
}

func BuildFolderMap(folders []types.Folder) map[string]string {
	folderMap := make(map[string]string)
	for _, folder := range folders {
		folderMap[folder.ID] = buildFolderPath(folder, folders)
	}
	return folderMap
}

func buildFolderPath(folder types.Folder, allFolders []types.Folder) string {
	path := folder.Name
	currentID := folder.ParentID

	for currentID != "" && currentID != "null" {
		parent := findFolderByID(currentID, allFolders)
		if parent == nil {
			break
		}
		path = fmt.Sprintf("%s/%s", parent.Name, path)
		currentID = parent.ParentID
	}

	return path
}

func findFolderByID(id string, folders []types.Folder) *types.Folder {
	for _, folder := range folders {
		if folder.ID == id {
			return &folder
		}
	}
	return nil
}

func ProcessSnippets(snippets []types.Snippet, folderMap map[string]string, tags []types.Tag) *types.OutputData {
	outputData := &types.OutputData{
		Snippets: make([]types.SnippetOutput, 0, len(snippets)),
	}

	// Create a map of tag IDs to tag names
	tagMap := make(map[string]string)
	for _, tag := range tags {
		tagMap[tag.ID] = tag.Name
	}

	for _, snippet := range snippets {
		snippetOutput := types.SnippetOutput{
			Name:       snippet.Name,
			FolderPath: getFolderPath(snippet.FolderID, folderMap),
			Tags:       make([]string, len(snippet.TagsIDs)),
			Content:    make([]types.ContentOutput, len(snippet.Content)),
		}

		// Convert tag IDs to tag names
		for i, tagID := range snippet.TagsIDs {
			if tagName, ok := tagMap[tagID]; ok {
				snippetOutput.Tags[i] = tagName
			} else {
				snippetOutput.Tags[i] = "Unknown Tag"
			}
		}

		for i, content := range snippet.Content {
			snippetOutput.Content[i] = types.ContentOutput{
				Label:    content.Label,
				Language: content.Language,
				Value:    content.Value,
			}
		}
		outputData.Snippets = append(outputData.Snippets, snippetOutput)
	}
	return outputData
}

func getFolderPath(folderID string, folderMap map[string]string) string {
	folderPath := folderMap[folderID]
	if folderPath == "" {
		folderPath = "Un-categorized"
	}
	return folderPath
}

func GenerateTextOutput(data *types.OutputData) (string, error) {
	var output strings.Builder
	for _, snippet := range data.Snippets {
		output.WriteString(fmt.Sprintf("Snippet: %s\n", snippet.Name))
		output.WriteString(fmt.Sprintf("Folder: %s\n", snippet.FolderPath))

		if len(snippet.Tags) > 0 {
			output.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(snippet.Tags, ", ")))
		}

		for _, content := range snippet.Content {
			output.WriteString(fmt.Sprintf("--- %s (%s) ---\n", content.Label, content.Language))
			output.WriteString(fmt.Sprintf("%s\n", content.Value))
			output.WriteString("\n")
		}

		output.WriteString(fmt.Sprintf("%s\n\n", strings.Repeat("-", 40)))
	}
	return output.String(), nil
}

func GenerateHTMLOutput(data *types.OutputData) (string, error) {
	// Convert the snippets to JSON
	snippetsJSON, err := json.Marshal(data.Snippets)
	if err != nil {
		return "", fmt.Errorf("error marshaling snippets to JSON: %v", err)
	}

	tmplContent, err := templates.TemplateFS.ReadFile("html/template.html")
	if err != nil {
		return "", fmt.Errorf("error reading embedded HTML template file: %v", err)
	}

	// Create a template with a custom function to safely inject JSON
	tmpl, err := template.New("output").Funcs(template.FuncMap{
		"safeJSON": func(v interface{}) template.JS {
			a, _ := json.Marshal(v)
			return template.JS(a)
		},
	}).Parse(string(tmplContent))
	if err != nil {
		return "", fmt.Errorf("error parsing HTML template: %v", err)
	}

	// Prepare the data for the template
	templateData := struct {
		SnippetsJSON template.JS
	}{
		SnippetsJSON: template.JS(snippetsJSON),
	}

	var output strings.Builder
	err = tmpl.Execute(&output, templateData)
	if err != nil {
		return "", fmt.Errorf("error executing HTML template: %v", err)
	}

	return output.String(), nil
}

func GenerateJSONOutput(data *types.OutputData) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}
	return string(jsonData), nil
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func getFileExtension(outputType string) string {
	switch outputType {
	case "html":
		return "html"
	case "json":
		return "json"
	default:
		return "txt"
	}
}
