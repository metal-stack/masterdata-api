package datastore

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"sigs.k8s.io/yaml"
)

type bootstrap[E Entity] struct {
	log *slog.Logger
	ds  Storage[E]
}

// Initdb reads all yaml files in given directory and apply their content as initial datasets.
func Initdb(log *slog.Logger, db *sqlx.DB, healthServer *health.Server, dir string) error {
	start := time.Now()
	files, err := filepath.Glob(path.Join(dir, "*.yaml"))
	if err != nil {
		return err
	}

	ts := New(log, db, &v1.Tenant{})
	tbs := &bootstrap[*v1.Tenant]{
		log: log,
		ds:  ts,
	}

	ps := New(log, db, &v1.Project{})
	pbs := &bootstrap[*v1.Project]{
		log: log,
		ds:  ps,
	}
	for _, f := range files {
		log.Info("read initdb for tenants", "file", f)
		err = tbs.processConfig(f)
		if err != nil {
			return err
		}
	}
	for _, f := range files {
		log.Info("read initdb for projects", "file", f)
		err = pbs.processConfig(f)
		if err != nil {
			return err
		}
	}
	log.Info("done reading initdb files", "took", time.Since(start))
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
func (ds *bootstrap[E]) processConfig(file string) error {
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

func (bs *bootstrap[E]) createOrUpdate(ctx context.Context, ydoc []byte) error {

	// all entities must contain a meta, parse that to get kind and apiversion
	var mm MetaMeta
	err := yaml.Unmarshal(ydoc, &mm)
	if err != nil {
		return err
	}
	bs.log.Info("initdb", "meta", mm.Meta.GetKind())

	kind := mm.Meta.GetKind()
	apiversion := mm.Meta.GetApiversion()

	var e E
	if kind != e.Kind() {
		bs.log.Info("skip", "kind from yaml", kind, "required kind", e.Kind())
		return nil
	}

	// messy extraction of apiversion from type
	if e.APIVersion() != apiversion {
		bs.log.Error("initdb apiversion does not match", "given", apiversion, "expected", e.APIVersion())
		return nil
	}

	ee := new(E)
	err = yaml.Unmarshal(ydoc, ee)
	if err != nil {
		return err
	}
	e = *ee

	newKind := e.GetMeta().GetKind()
	newID := e.GetMeta().GetId()

	// now check that this type is already present for this id,
	// therefore create nil interface to get into
	exists := true
	existingEntity, err := bs.ds.Get(ctx, mm.Meta.GetId())
	if err != nil {
		if errors.As(err, &NotFoundError{}) {
			exists = false
		} else {
			bs.log.Error("initdb", "error", err)
			return err
		}
	}
	// now check if it exists by checking for id presence
	// then update, otherwise create
	if exists {
		oldVersion := existingEntity.GetMeta().GetVersion()
		newVersion := e.GetMeta().GetVersion()
		bs.log.Info("initdb found existing, update", "kind", newKind, "id", newID, "old version", oldVersion, "new version", newVersion)
		if oldVersion >= newVersion {
			bs.log.Info("initdb existing version is equal or higher, skip update", "kind", newKind, "id", newID)
			return nil
		}

		e.GetMeta().SetVersion(existingEntity.GetMeta().GetVersion())
		err = bs.ds.Update(ctx, e)
		if err != nil {
			return err
		}
	} else {
		bs.log.Info("initdb create", "kind", newKind, "id", newID)
		err = bs.ds.Create(ctx, e)
		if err != nil {
			return err
		}
	}
	return nil
}
