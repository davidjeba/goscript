package workspace

import "fmt"

// BlockKind identifies a document block type.
type BlockKind string

const (
	BlockPage      BlockKind = "page"
	BlockHeading   BlockKind = "heading"
	BlockParagraph BlockKind = "paragraph"
	BlockBulleted  BlockKind = "bulleted_list"
	BlockNumbered  BlockKind = "numbered_list"
	BlockTodo      BlockKind = "todo"
	BlockToggle    BlockKind = "toggle"
	BlockCode      BlockKind = "code"
	BlockQuote     BlockKind = "quote"
	BlockCallout   BlockKind = "callout"
	BlockTable     BlockKind = "table"
	BlockDatabase  BlockKind = "database"
	BlockFile      BlockKind = "file"
)

// Block is a Notion-style document node.
type Block struct {
	ID       string                 `json:"id"`
	ParentID string                 `json:"parentId,omitempty"`
	Kind     BlockKind              `json:"kind"`
	Text     string                 `json:"text,omitempty"`
	Props    map[string]interface{} `json:"props,omitempty"`
	Children []*Block               `json:"children,omitempty"`
}

// Document is a block-based document tree.
type Document struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Blocks []*Block `json:"blocks"`
}

// AddBlock adds a block at the document root.
func (d *Document) AddBlock(block *Block) {
	if block == nil {
		return
	}
	d.Blocks = append(d.Blocks, block)
}

// FindBlock finds a block by ID anywhere in the tree.
func (d *Document) FindBlock(id string) (*Block, bool) {
	for _, block := range d.Blocks {
		if found, ok := findBlock(block, id); ok {
			return found, true
		}
	}
	return nil, false
}

// AppendChild adds a child block to a parent block.
func (d *Document) AppendChild(parentID string, child *Block) error {
	parent, ok := d.FindBlock(parentID)
	if !ok {
		return fmt.Errorf("parent block %q not found", parentID)
	}

	child.ParentID = parentID
	parent.Children = append(parent.Children, child)
	return nil
}

// Flatten returns all blocks in depth-first order.
func (d *Document) Flatten() []*Block {
	var flattened []*Block
	for _, block := range d.Blocks {
		flattened = append(flattened, flattenBlock(block)...)
	}
	return flattened
}

func findBlock(block *Block, id string) (*Block, bool) {
	if block == nil {
		return nil, false
	}
	if block.ID == id {
		return block, true
	}
	for _, child := range block.Children {
		if found, ok := findBlock(child, id); ok {
			return found, true
		}
	}
	return nil, false
}

func flattenBlock(block *Block) []*Block {
	if block == nil {
		return nil
	}

	result := []*Block{block}
	for _, child := range block.Children {
		result = append(result, flattenBlock(child)...)
	}
	return result
}
