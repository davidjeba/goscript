package workspace

import "testing"

func TestDocumentFindBlockAndFlatten(t *testing.T) {
	doc := &Document{ID: "doc-1", Title: "Specs"}
	root := &Block{ID: "root", Kind: BlockPage, Text: "Root"}
	child := &Block{ID: "child", Kind: BlockParagraph, Text: "Hello"}

	doc.AddBlock(root)
	if err := doc.AppendChild("root", child); err != nil {
		t.Fatalf("AppendChild returned error: %v", err)
	}

	found, ok := doc.FindBlock("child")
	if !ok || found.Text != "Hello" {
		t.Fatalf("expected to find child block")
	}

	if len(doc.Flatten()) != 2 {
		t.Fatalf("expected flattened length 2, got %d", len(doc.Flatten()))
	}
}

func TestTableUpdateCell(t *testing.T) {
	table := &Table{
		ID:   "table-1",
		Name: "Tasks",
		Columns: []TableColumn{
			{ID: "title", Name: "Title", Type: ColumnText, Editable: true},
			{ID: "done", Name: "Done", Type: ColumnBool, Editable: true},
		},
		Rows: []TableRow{
			{ID: "row-1", Values: map[string]interface{}{"title": "Ship Vibe"}},
		},
	}

	if err := table.UpdateCell("row-1", "done", true); err != nil {
		t.Fatalf("UpdateCell returned error: %v", err)
	}

	if table.Rows[0].Values["done"] != true {
		t.Fatalf("expected done=true, got %v", table.Rows[0].Values["done"])
	}
}

func TestFileNodeFindPath(t *testing.T) {
	root := &FileNode{Name: "root", Path: "/", Type: NodeDirectory}
	docs := &FileNode{Name: "docs", Path: "/docs", Type: NodeDirectory}
	file := &FileNode{Name: "spec.md", Path: "/docs/spec.md", Type: NodeFile}

	root.AddChild(docs)
	docs.AddChild(file)

	found, ok := root.FindPath("/docs/spec.md")
	if !ok || found.Name != "spec.md" {
		t.Fatalf("expected to find spec.md")
	}
}
