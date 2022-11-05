package datastore

import (
	"context"
	"errors"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	healthv1 "github.com/metal-stack/masterdata-api/api/grpc/health/v1"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/health"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// FIXME iterate over all E here

// Initdb reads all yaml files in given directory and apply their content as initial datasets.
func (ds *datastore[E]) Initdb(healthServer *health.Server, dir string) error {
	files, err := filepath.Glob(path.Join(dir, "*.yaml"))
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

	var e E
	healthServer.SetServingStatus("initdb-"+e.Kind(), healthv1.HealthCheckResponse_SERVING)
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
func (ds *datastore[E]) processConfig(file string) error {
	yml, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	ctx := context.Background()

	yamldocs := splitYamlDocs(string(yml))
	for i := range yamldocs {
		ydoc := yamldocs[i]
		err = ds.createOrUpdate(ctx, []byte(ydoc))
		if err != nil {
			return err
		}
	}
	return nil
}

func (ds *datastore[E]) createOrUpdate(ctx context.Context, ydoc []byte) error {

	// all entities must contain a meta, parse that to get kind and apiversion
	var mm MetaMeta
	err := yaml.Unmarshal(ydoc, &mm)
	if err != nil {
		return err
	}
	ds.log.Info("initdb", zap.Any("meta", mm.Meta.GetKind()))

	kind := mm.Meta.GetKind()
	apiversion := mm.Meta.GetApiversion()
	var e E
	if kind != e.Kind() {
		ds.log.Info("skip", zap.String("kind from yaml", kind), zap.String("required kind", e.Kind()))
		return nil
	}

	// messy extraction of apiversion from type
	if e.APIVersion() != apiversion {
		ds.log.Error("initdb apiversion does not match", zap.String("given", apiversion), zap.String("expected", e.APIVersion()))
		return nil
	}

	err = yaml.Unmarshal(ydoc, e)
	if err != nil {
		return err
	}

	newKind := e.GetMeta().GetKind()
	newID := e.GetMeta().GetId()

	// now check that this type is already present for this id,
	// therefore create nil interface to get into
	exists := true
	existingEntity, err := ds.Get(ctx, mm.Meta.GetId())
	if err != nil {
		if errors.As(err, &NotFoundError{}) {
			exists = false
		} else {
			ds.log.Error("initdb", zap.Error(err))
			return err
		}
	}
	// now check if it exists by checking for id presence
	// then update, otherwise create
	if exists {
		oldVersion := existingEntity.GetMeta().GetVersion()
		newVersion := e.GetMeta().GetVersion()
		ds.log.Info("initdb found existing, update", zap.String("kind", newKind), zap.String("id", newID), zap.Any("old version", oldVersion), zap.Any("new version", newVersion))
		if oldVersion >= newVersion {
			ds.log.Info("initdb existing version is equal or higher, skip update", zap.String("kind", newKind), zap.String("id", newID))
			return nil
		}

		e.GetMeta().SetVersion(existingEntity.GetMeta().GetVersion())
		err = ds.Update(ctx, e)
		if err != nil {
			return err
		}
	} else {
		ds.log.Info("initdb create", zap.String("kind", newKind), zap.String("id", newID))
		err = ds.Create(ctx, e)
		if err != nil {
			return err
		}
	}
	return nil
}
