package workspace

import "strings"

// FileNodeType identifies whether a node is a file or directory.
type FileNodeType string

const (
	NodeFile      FileNodeType = "file"
	NodeDirectory FileNodeType = "directory"
)

// FileNode models a file-browser tree.
type FileNode struct {
	Name     string            `json:"name"`
	Path     string            `json:"path"`
	Type     FileNodeType      `json:"type"`
	Children []*FileNode       `json:"children,omitempty"`
	Meta     map[string]string `json:"meta,omitempty"`
}

// AddChild adds a child node to a directory.
func (n *FileNode) AddChild(child *FileNode) {
	if n == nil || child == nil {
		return
	}
	n.Children = append(n.Children, child)
}

// FindPath finds a node by its logical path.
func (n *FileNode) FindPath(path string) (*FileNode, bool) {
	if n == nil {
		return nil, false
	}

	normalized := normalizePath(path)
	if normalizePath(n.Path) == normalized {
		return n, true
	}

	for _, child := range n.Children {
		if found, ok := child.FindPath(normalized); ok {
			return found, true
		}
	}

	return nil, false
}

func normalizePath(path string) string {
	path = strings.ReplaceAll(strings.TrimSpace(path), "\\", "/")
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}
