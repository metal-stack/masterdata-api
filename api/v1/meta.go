package apiv1

import (
	jsoniter "github.com/json-iterator/go"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (m *Meta) SetId(id string) {
	m.Id = id
}

func (m *Meta) SetVersion(version int64) {
	m.Version = version
}

func (m *Meta) SetCreatedTime(time *timestamppb.Timestamp) {
	m.CreatedTime = time
}

func (m *Meta) SetUpdatedTime(time *timestamppb.Timestamp) {
	m.UpdatedTime = time
}

func (m *Meta) SetAnnotations(annotations map[string]string) {
	m.Annotations = annotations
}

func (m *Meta) SetLabels(labels []string) {
	m.Labels = labels
}
