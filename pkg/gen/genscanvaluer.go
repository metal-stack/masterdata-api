// Generator should not be part of the package
//+build ignore

/*
	Generates methods necessary to fulfil the Scanner and Valuer Interface.

	Example Usage for Type "Tenant" in package "v1"
	//go:generate go run ../../pkg/gen/gensv.go -package v1 -type Tenant
*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	errs "github.com/pkg/errors"
	"go.uber.org/zap"
)

const defaultFilenameSuffix = "_scnrvalr.go"

var (
	packageName = flag.String("package", "", "package name; must be set")
	typeName    = flag.String("type", "", "type name; must be set")
	output      = flag.String("output", "", "output file name; default srcdir/<type>_scnrvalr.go")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of gensv:\n")
	fmt.Fprintf(os.Stderr, "\tgensv -package P -type T\n")
	fmt.Fprintf(os.Stderr, "\tgensv -package [package] -type [type] -o [filename]\n")
}

func main() {
	log.SetPrefix("gensv: ")
	flag.Usage = Usage
	flag.Parse()

	if len(*packageName) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	if len(*typeName) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// Parse the package once.
	var dir string
	g := Generator{}

	err := g.generate(*packageName, *typeName)
	if err != nil {
		log.Fatal("error generating content", zap.Error(err))
	}

	// Format the output.
	src := g.format()

	// Write to file.
	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s%s", *typeName, defaultFilenameSuffix)
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}
	err = ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatal("error writing output", zap.Error(err))
	}
}

// Generator is primarily used to buffer the output for format.Source.
type Generator struct {
	buf bytes.Buffer // Accumulated output.
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// generate produces the content of the file for the given type.
func (g *Generator) generate(packageName, typeName string) error {

	tmpl, err := template.New("sv").Parse(svTemplate)
	if err != nil {
		return errs.Wrapf(err, "error parsing template %s", svTemplate)
	}
	stmpl, err := template.New("sv").Parse(schemaTemplate)
	if err != nil {
		return errs.Wrapf(err, "error parsing template %s", schemaTemplate)
	}

	info := map[string]string{
		"packageName":   packageName,
		"typeName":      typeName,
		"typeNameLower": strings.ToLower(typeName),
		"tableName":     strings.ToLower(typeName) + "s",
	}

	var renderedBytesSchema bytes.Buffer
	err = stmpl.Execute(&renderedBytesSchema, info)
	if err != nil {
		return errs.Wrap(err, "error rendering template")
	}

	info["schema"] = fmt.Sprintf("`%s`", renderedBytesSchema.String())

	var renderedBytes bytes.Buffer
	err = tmpl.Execute(&renderedBytes, info)
	if err != nil {
		return errs.Wrap(err, "error rendering template")
	}

	g.Printf("%s", renderedBytes.String())
	return nil
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}

const schemaTemplate = `
	CREATE TABLE IF NOT EXISTS {{ .tableName }} (
		id   text PRIMARY KEY NOT NULL,
		{{ .typeNameLower }} JSONB NOT NULL
	);
	CREATE INDEX IF NOT EXISTS {{ .typeNameLower }}_idx ON {{ .tableName }} USING GIN({{ .typeNameLower }});

	CREATE TABLE IF NOT EXISTS {{ .tableName }}_history (
		id         text NOT NULL,
		op		   char NOT NULL,
		created_at timestamptz NOT NULL,
		{{ .typeNameLower }} JSONB NOT NULL
	);
	CREATE INDEX IF NOT EXISTS {{ .tableName }}_history_id_created_at_idx ON {{ .tableName }}_history USING btree(id, created_at);
`

const svTemplate = `
// This file was automatically generated by pkg/gen/genscanvaluer.
// DO NOT EDIT MANUALLY.
// Regenerate with "go generate" or "make generate"

package {{ .packageName }}

import (
	"database/sql/driver"
	"fmt"
)

func (m {{ .typeName }}) Schema() string {
	return {{ .schema }}
}

func (m {{ .typeName }}) JSONField() string {
	return "{{ .typeNameLower }}"
}

func (m {{ .typeName }}) TableName() string {
	return "{{ .typeNameLower }}s"
}

func (m {{ .typeName }}) Kind() string {
	return "{{ .typeName }}"
}

func (m {{ .typeName }}) APIVersion() string {
	return "{{ .packageName }}"
}

// Value make the {{ .typeName }} struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (m {{ .typeName }}) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan make the {{ .typeName }} struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (m *{{ .typeName }}) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}

	err := json.Unmarshal(b, m)
	return err
}`
