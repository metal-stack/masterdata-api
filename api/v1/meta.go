package v1

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (m *Meta) SetId(id string) {
	m.Id = id
}

func (m *Meta) SetVersion(version int64) {
	m.Version = version
}

func (m *Meta) SetCreatedTime(time *timestamp.Timestamp) {
	m.CreatedTime = time
}

func (m *Meta) SetUpdatedTime(time *timestamp.Timestamp) {
	m.UpdatedTime = time
}
