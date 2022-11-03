package v1

import (
	"errors"
	"reflect"
	"time"
)

type QuotaSet struct {
	Cluster *Quota `json:"cluster,omitempty"`
	Machine *Quota `json:"machine,omitempty"`
	Ip      *Quota `json:"ip,omitempty"`
	Project *Quota `json:"project,omitempty"`
}

type Quota struct {
	Used  *int32 `json:"used,omitempty"`
	Quota *int32 `json:"quota,omitempty"`
}

type Meta struct {
	Id          string            `json:"id,omitempty"`
	Kind        string            `json:"kind,omitempty"`
	Apiversion  string            `json:"apiversion,omitempty"`
	Version     int64             `json:"version,omitempty"`
	CreatedTime *time.Time        `json:"created_time,omitempty"`
	UpdatedTime *time.Time        `json:"updated_time,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      []string          `json:"labels,omitempty"`
}

type (
	ProjectIDTime struct {
		ProjectID string    `json:"project_id,omitempty" description:"projectID as returned by cloud-api (e.g. 10241dd7-a8de-4856-8ac0-b55830b22036)"`
		Time      time.Time `json:"time,omitempty" description:"point in time"`
	}

	// cluster-identification e.g. from prometheus-metrics
	ClusterNameProject struct {
		ClusterName string `json:"cluster_name,omitempty" description:"cluster name"`
		Project     string `json:"project,omitempty" description:"generated middle-part of gardener shoot namespace, e.g. 'ps5d42'"`
	}

	// MasterdataLookupRequest to specify which masterdata tenant and project to return
	MasterdataLookupRequest struct {
		// search by name and project
		ClusterNameProject *ClusterNameProject `json:"cluster_name_project,omitempty" description:"lookup by clustername and shoot-project"`

		// OR search by ClusterID as returned by cloud-api (e.g. 345abc12-3321-4dbc-8d17-55c6ea4fcb38)
		ClusterID *string `json:"cluster_id,omitempty" description:"lookup by clusterID as returned by cloud-api (e.g. 345abc12-3321-4dbc-8d17-55c6ea4fcb38)"`

		// OR search by ProjectID as returned by cloud-api (e.g. 10241dd7-a8de-4856-8ac0-b55830b22036) at some point in time
		ProjectIDTime *ProjectIDTime `json:"project_id_time,omitempty" description:"lookup at some point in time by projectID as returned by cloud-api (e.g. 10241dd7-a8de-4856-8ac0-b55830b22036)"`
	}

	// MasterdataLookupResponse contains the masterdata tenant and project
	MasterdataLookupResponse struct {
		Tenant  *Tenant  `json:"tenant,omitempty" description:"tenant to which the project belongs"`
		Project *Project `json:"project,omitempty" description:"project"`
	}
)

// Validate validates the request.
func (mlr *MasterdataLookupRequest) Validate() error {
	if mlr == nil {
		return errors.New("masterdataLookupRequest is nil")
	}
	if mlr.ClusterID == nil && mlr.ClusterNameProject == nil && mlr.ProjectIDTime == nil {
		return errors.New("one of ClusterNameProject, ClusterID or ProjectID must be set")
	}
	if !onlyOneOfPtrsSet(mlr.ClusterID, mlr.ClusterNameProject, mlr.ProjectIDTime) {
		return errors.New("only one of ClusterNameProject, ClusterID or ProjectIDTime may be set")
	}

	return nil
}

// onlyOneOfPtrsSet returns true if only one of the given pointers is not nil.
// If values are given they are always counted as not nil.
func onlyOneOfPtrsSet(ptrs ...any) bool {
	count := 0
	for _, p := range ptrs {
		rv := reflect.ValueOf(p)
		if rv.Kind() != reflect.Ptr || !rv.IsNil() {
			count++
		}
	}
	return count == 1
}
