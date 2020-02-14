package datastore

import (
	"context"
	"fmt"
	gyaml "github.com/ghodss/yaml"
	"io/ioutil"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	healthv1 "github.com/metal-stack/masterdata-api/api/grpc/health/v1"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/health"
)

// Initdb reads all yaml files in given directory and apply their content as initial datasets.
func (d *Datastore) Initdb(healthServer *health.Server, dir string) error {
	files, err := d.listFiles(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		d.log.Sugar().Infow("read initdb", "file", f)
		err = d.processConfig(f)
		if err != nil {
			return err
		}
	}
	healthServer.SetServingStatus("initdb", healthv1.HealthCheckResponse_SERVING)
	return nil
}

// MetaMeta is a container for the meta property inside a entity.
type MetaMeta struct {
	v1.Meta `json:"meta" yaml:"meta"`
}

// activate multiline-mode so that ^ matches start of line
var docSplitExpr = regexp.MustCompile(`(?m)^---\s*\n`)

// splitYamlDocs splits the given (possible multi-doc) yamldoc in single documents, skips empty docs.
// If doc is blank, nil is returned.
func splitYamlDocs(doc string) []string {

	docs := docSplitExpr.Split(doc, -1)
	var result []string
	for i := range docs {
		// only append non-empty docs
		if len(strings.TrimSpace(docs[i])) > 0 {
			result = append(result, docs[i])
		}
	}
	return result
}

// processConfig processes all yaml docs contained in the given file
func (d *Datastore) processConfig(file string) error {
	yaml, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	ctx := context.Background()

	yamldocs := splitYamlDocs(string(yaml))
	for i := range yamldocs {
		ydoc := yamldocs[i]
		err = d.createOrUpdate(ctx, []byte(ydoc))
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Datastore) createOrUpdate(ctx context.Context, ydoc []byte) error {

	// all entities must contain a meta, parse that to get kind and apiversion
	var mm MetaMeta
	err := gyaml.Unmarshal(ydoc, &mm)
	if err != nil {
		return err
	}
	meta := mm.Meta
	d.log.Sugar().Infow("initdb", "meta", meta)

	kind := meta.GetKind()
	apiversion := meta.GetApiversion()

	// get type for this kind from datastore entity types registry
	typeElem, ok := d.types[strings.ToLower(kind)]
	if !ok {
		return fmt.Errorf("initdb: unknown kind:%s", kind)
	}
	// messy extraction of apiversion from type
	// type is like "*v1.Project"
	elementType := reflect.TypeOf(typeElem)
	expectedAPI := strings.ReplaceAll(strings.Split(elementType.String(), ".")[0], "*", "")
	if expectedAPI != apiversion {
		d.log.Sugar().Errorw("initdb apiversion does not match", "given", apiversion, "expected", expectedAPI)
		return nil
	}
	d.log.Sugar().Infow("initdb", "type", elementType, "apiversion", apiversion)

	// no we have the type, create new from type and marshall in that new struct
	newEntity, ok := reflect.New(elementType.Elem()).Interface().(VersionedJSONEntity)
	if !ok {
		panic(fmt.Sprintf("entity type %s must implement VersionedJSONEntity-Interface", elementType.String()))
	}

	err = gyaml.Unmarshal(ydoc, newEntity)
	if err != nil {
		return err
	}

	newKind := newEntity.GetMeta().GetKind()
	newID := newEntity.GetMeta().GetId()

	// now check that this type is already present for this id,
	// therefore create nil interface to get into
	existingEntity := reflect.New(elementType.Elem()).Interface().(VersionedJSONEntity)
	err = d.Get(ctx, meta.GetId(), existingEntity)
	if err != nil {
		switch err.(type) {
		case NotFoundError:
			existingEntity = nil
		default:
			d.log.Sugar().Errorw("initdb", "error", err)
			return err
		}
	}
	// now check if it exists by checking for id presence
	// then update, otherwise create
	if existingEntity != nil {
		oldVersion := existingEntity.GetMeta().GetVersion()
		newVersion := newEntity.GetMeta().GetVersion()
		d.log.Sugar().Infow("initdb found existing, update", "kind", newKind, "id", newID, "old version", oldVersion, "new version", newVersion)
		if oldVersion >= newVersion {
			d.log.Sugar().Infow("initdb existing version is equal or higher, skip update", "kind", newKind, "id", newID)
			return nil
		}

		newEntity.GetMeta().SetVersion(existingEntity.GetMeta().GetVersion())
		err = d.Update(ctx, newEntity)
		if err != nil {
			return err
		}
	} else {
		d.log.Sugar().Infow("initdb create", "kind", newKind, "id", newID)
		err = d.Create(context.Background(), newEntity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Datastore) listFiles(dir string) ([]string, error) {
	matches, err := filepath.Glob(path.Join(dir, "*.yaml"))
	if err != nil {
		return nil, err
	}
	return matches, nil
}
