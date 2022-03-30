// This file was automatically generated by pkg/gen/genscanvaluer.
// DO NOT EDIT MANUALLY.
// Regenerate with "go generate" or "make generate"

package v1

import (
	"database/sql/driver"
	"fmt"
)

func (p *Person) Schema() string {
	return `
	CREATE TABLE IF NOT EXISTS persons (
		id   text PRIMARY KEY NOT NULL,
		person JSONB NOT NULL
	);
	CREATE INDEX IF NOT EXISTS person_idx ON persons USING GIN(person);

	CREATE TABLE IF NOT EXISTS persons_history (
		id         text NOT NULL,
		op         char NOT NULL,
		created_at timestamptz NOT NULL,
		person JSONB NOT NULL
	);
	CREATE INDEX IF NOT EXISTS persons_history_id_created_at_idx ON persons_history USING btree(id, created_at);
`
}

func (p *Person) JSONField() string {
	return "person"
}

func (p *Person) TableName() string {
	return "persons"
}

func (p *Person) Kind() string {
	return "Person"
}

func (p *Person) APIVersion() string {
	return "v1"
}

// Value make the Person struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (p *Person) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan make the Person struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (p *Person) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}

	err := json.Unmarshal(b, p)
	return err
}
