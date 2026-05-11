package workspace

import "fmt"

// ColumnType identifies a database/table field type.
type ColumnType string

const (
	ColumnText        ColumnType = "text"
	ColumnNumber      ColumnType = "number"
	ColumnBool        ColumnType = "bool"
	ColumnDate        ColumnType = "date"
	ColumnSelect      ColumnType = "select"
	ColumnMultiSelect ColumnType = "multi_select"
	ColumnRelation    ColumnType = "relation"
	ColumnFile        ColumnType = "file"
)

// TableColumn describes a schema field.
type TableColumn struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Type     ColumnType `json:"type"`
	Editable bool       `json:"editable"`
	Width    int        `json:"width,omitempty"`
}

// TableRow holds cell values.
type TableRow struct {
	ID     string                 `json:"id"`
	Values map[string]interface{} `json:"values"`
}

// Table models an inline-editable database table.
type Table struct {
	ID      string        `json:"id"`
	Name    string        `json:"name"`
	Columns []TableColumn `json:"columns"`
	Rows    []TableRow    `json:"rows"`
}

// AddColumn adds a column to the schema.
func (t *Table) AddColumn(column TableColumn) {
	t.Columns = append(t.Columns, column)
}

// AddRow appends a row to the table.
func (t *Table) AddRow(row TableRow) {
	if row.Values == nil {
		row.Values = map[string]interface{}{}
	}
	t.Rows = append(t.Rows, row)
}

// UpdateCell changes one cell value after schema validation.
func (t *Table) UpdateCell(rowID, columnID string, value interface{}) error {
	column, ok := t.columnByID(columnID)
	if !ok {
		return fmt.Errorf("column %q not found", columnID)
	}
	if !column.Editable {
		return fmt.Errorf("column %q is not editable", columnID)
	}

	for rowIndex := range t.Rows {
		if t.Rows[rowIndex].ID != rowID {
			continue
		}
		if t.Rows[rowIndex].Values == nil {
			t.Rows[rowIndex].Values = map[string]interface{}{}
		}
		t.Rows[rowIndex].Values[columnID] = value
		return nil
	}

	return fmt.Errorf("row %q not found", rowID)
}

func (t *Table) columnByID(columnID string) (TableColumn, bool) {
	for _, column := range t.Columns {
		if column.ID == columnID {
			return column, true
		}
	}
	return TableColumn{}, false
}
