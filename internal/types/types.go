package types

type Database struct {
	Folders  []Folder  `json:"folders"`
	Snippets []Snippet `json:"snippets"`
	Tags     []Tag     `json:"tags"`
}

type Folder struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	ParentID        string `json:"parentId"`
	IsOpen          bool   `json:"isOpen"`
	DefaultLanguage string `json:"defaultLanguage"`
	Index           int    `json:"index"`
	CreatedAt       int64  `json:"createdAt"`
	UpdatedAt       int64  `json:"updatedAt"`
}

type Snippet struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	FolderID    string          `json:"folderId"`
	TagsIDs     []string        `json:"tagsIds"`
	Content     []ContentOutput `json:"content"`
	IsFavorites bool            `json:"isFavorites"`
	IsDeleted   bool            `json:"isDeleted"`
	CreatedAt   int64           `json:"createdAt"`
	UpdatedAt   int64           `json:"updatedAt"`
}

type Tag struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

type OutputData struct {
	Snippets []SnippetOutput `json:"snippets"`
}

type SnippetOutput struct {
	Name       string          `json:"name"`
	FolderPath string          `json:"folderPath"`
	Tags       []string        `json:"tags"`
	Content    []ContentOutput `json:"content"`
}

type ContentOutput struct {
	Label    string `json:"label"`
	Language string `json:"language"`
	Value    string `json:"value"`
}
