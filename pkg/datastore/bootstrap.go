package datastore

import (
	"context"
	"fmt"
	gyaml "github.com/ghodss/yaml"
	"github.com/golang/protobuf/ptypes"
	healthv1 "github.com/metal-stack/masterdata-api/api/grpc/health/v1"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/health"
	"go.uber.org/zap"
	"io/ioutil"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

// Initdb reads all yaml files in given directory and apply their content as initial datasets.
func (ds *Datastore) Initdb(healthServer *health.Server, dir string) error {
	files, err := ds.listFiles(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		ds.log.Info("read initdb", zap.Any("file", f))
		err = ds.processConfig(f)
		if err != nil {
			return err
		}
	}

	entities := []interface{}{&[]v1.Project{}, &[]v1.Tenant{}}
	for _, e := range entities {
		err = ds.consolidateHistory(e)
		if err != nil {
			ds.log.Error("error consolidate history", zap.Error(err))
		}
	}

	healthServer.SetServingStatus("initdb", healthv1.HealthCheckResponse_SERVING)
	return nil
}

// consolidateHistory ensures, that for each VersionedJSONEntity there is at least one "created"-row in the history table.
// The type of entities to consolidate is specified by the given pointer to a slice of entities.
func (ds *Datastore) consolidateHistory(entitySlicePrt interface{}) error {

	entitySliceV := reflect.ValueOf(entitySlicePrt)
	if entitySliceV.Kind() != reflect.Ptr || entitySliceV.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("entity argument must be a slice address")
	}
	entitySliceElem := entitySliceV.Elem()
	entitySliceElementType := entitySliceElem.Type().Elem()

	tx, err := ds.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer ds.rollback(tx)

	filter := make(map[string]interface{})
	err = ds.Find(context.Background(), filter, entitySlicePrt)
	if err != nil {
		return err
	}

	for i := 0; i < entitySliceElem.Len(); i++ {
		vpi := entitySliceElem.Index(i).Addr().Interface()
		enityVe, ok := vpi.(VersionedJSONEntity)
		if !ok {
			return fmt.Errorf("element type must implement VersionedJSONEntity-Interface, was %T", vpi)
		}

		historyVe, ok := reflect.New(entitySliceElementType).Interface().(VersionedJSONEntity)
		if !ok {
			return fmt.Errorf("element type must implement VersionedJSONEntity-Interface")
		}

		// check if we already have a "created" row for this entity in history
		err = ds.GetHistoryCreated(context.Background(), enityVe.GetMeta().Id, historyVe)
		if err == nil {
			continue
		}
		if _, notFound := err.(NotFoundError); !notFound {
			return err // some sort of technical error stops us
		}

		// consolidate history by inserting the "created" row in history at the correct date
		entityCreatedPbTs := enityVe.GetMeta().CreatedTime
		entityCreated, err := ptypes.Timestamp(entityCreatedPbTs)
		if err != nil {
			return err
		}
		err = ds.insertHistory(enityVe, opCreate, entityCreated, tx)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
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
func (ds *Datastore) processConfig(file string) error {
	yaml, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	ctx := context.Background()

	yamldocs := splitYamlDocs(string(yaml))
	for i := range yamldocs {
		ydoc := yamldocs[i]
		err = ds.createOrUpdate(ctx, []byte(ydoc))
		if err != nil {
			return err
		}
	}
	return nil
}

func (ds *Datastore) createOrUpdate(ctx context.Context, ydoc []byte) error {

	// all entities must contain a meta, parse that to get kind and apiversion
	var mm MetaMeta
	err := gyaml.Unmarshal(ydoc, &mm)
	if err != nil {
		return err
	}
	meta := mm.Meta
	ds.log.Info("initdb", zap.Any("meta", meta))

	kind := meta.GetKind()
	apiversion := meta.GetApiversion()

	// get type for this kind from datastore entity types registry
	typeElem, ok := ds.types[strings.ToLower(kind)]
	if !ok {
		return fmt.Errorf("initdb: unknown kind:%s", kind)
	}
	// messy extraction of apiversion from type
	// type is like "*v1.Project"
	elementType := reflect.TypeOf(typeElem)
	expectedAPI := strings.ReplaceAll(strings.Split(elementType.String(), ".")[0], "*", "")
	if expectedAPI != apiversion {
		ds.log.Error("initdb apiversion does not match", zap.String("given", apiversion), zap.String("expected", expectedAPI))
		return nil
	}
	ds.log.Info("initdb", zap.Stringer("type", elementType), zap.String("apiversion", apiversion))

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
	err = ds.Get(ctx, meta.GetId(), existingEntity)
	if err != nil {
		switch err.(type) {
		case NotFoundError:
			existingEntity = nil
		default:
			ds.log.Error("initdb", zap.Error(err))
			return err
		}
	}
	// now check if it exists by checking for id presence
	// then update, otherwise create
	if existingEntity != nil {
		oldVersion := existingEntity.GetMeta().GetVersion()
		newVersion := newEntity.GetMeta().GetVersion()
		ds.log.Info("initdb found existing, update", zap.String("kind", newKind), zap.String("id", newID), zap.Any("old version", oldVersion), zap.Any("new version", newVersion))
		if oldVersion >= newVersion {
			ds.log.Info("initdb existing version is equal or higher, skip update", zap.String("kind", newKind), zap.String("id", newID))
			return nil
		}

		newEntity.GetMeta().SetVersion(existingEntity.GetMeta().GetVersion())
		err = ds.Update(ctx, newEntity)
		if err != nil {
			return err
		}
	} else {
		ds.log.Info("initdb create", zap.String("kind", newKind), zap.String("id", newID))
		err = ds.Create(ctx, newEntity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ds *Datastore) listFiles(dir string) ([]string, error) {
	matches, err := filepath.Glob(path.Join(dir, "*.yaml"))
	if err != nil {
		return nil, err
	}
	return matches, nil
}
