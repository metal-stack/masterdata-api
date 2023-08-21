// This file was automatically generated by pkg/gen/genscanvaluer.
// DO NOT EDIT MANUALLY.
// Regenerate with "go generate" or "make generate"

package v1

import (
	"database/sql/driver"
	"fmt"
)

func (p *Project) Schema() string {
	return `
	CREATE TABLE IF NOT EXISTS projects (
		id   text PRIMARY KEY NOT NULL,
		project JSONB NOT NULL
	);
	CREATE INDEX IF NOT EXISTS project_idx ON projects USING GIN(project);

	CREATE TABLE IF NOT EXISTS projects_history (
		id         text NOT NULL,
		op         char NOT NULL,
		created_at timestamptz NOT NULL,
		project JSONB NOT NULL
	);
	CREATE INDEX IF NOT EXISTS projects_history_id_created_at_idx ON projects_history USING btree(id, created_at);
`
}

func (p *Project) JSONField() string {
	return "project"
}

func (p *Project) TableName() string {
	return "projects"
}

func (p *Project) Kind() string {
	return "Project"
}

func (p *Project) APIVersion() string {
	return "v1"
}

// Value make the Project struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (p *Project) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan make the Project struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (p *Project) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}

	err := json.Unmarshal(b, p)
	return err
}
