// This file was automatically generated by pkg/gen/genscanvaluer.
// DO NOT EDIT MANUALLY.
// Regenerate with "go generate" or "make generate"

package v1

import (
	"database/sql/driver"
	"fmt"
)

func (m Project) Schema() string {
	return `
	CREATE TABLE IF NOT EXISTS projects (
		id   text PRIMARY KEY NOT NULL,
		project JSONB NOT NULL
	);
	CREATE INDEX IF NOT EXISTS project_idx ON projects USING GIN(project);
`
}

func (m Project) JSONField() string {
	return "project"
}

func (m Project) TableName() string {
	return "projects"
}

// Value make the Project struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (m Project) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan make the Project struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (m *Project) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}

	err := json.Unmarshal(b, m)
	return err
}
