package mapper

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	v1 "github.com/metal-stack/masterdata-api/api/rest/v1"
	mdmv1 "github.com/metal-stack/masterdata-api/api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ToMdmV1Tenant(t *v1.Tenant) *mdmv1.Tenant {
	if t == nil {
		return nil
	}

	return &mdmv1.Tenant{
		Meta:          ToMdmV1Meta(t.Meta),
		Name:          t.Name,
		Description:   t.Description,
		Quotas:        ToMdmV1QuotaSet(t.Quotas),
		DefaultQuotas: ToMdmV1QuotaSet(t.DefaultQuotas),
		IamConfig:     ToMdmV1IamConfig(t.IAMConfig),
	}
}

func ToMdmV1IamConfig(i *v1.IAMConfig) *mdmv1.IAMConfig {
	if i == nil {
		return nil
	}
	return &mdmv1.IAMConfig{
		IssuerConfig: ToMdmV1IssuerConfig(i.IssuerConfig),
		IdmConfig:    ToMdmV1IDMConfig(i.IDMConfig),
	}
}

func ToMdmV1IssuerConfig(i *v1.IssuerConfig) *mdmv1.IssuerConfig {
	if i == nil {
		return nil
	}
	return &mdmv1.IssuerConfig{
		Url:      i.URL,
		ClientId: i.ClientID,
	}
}

func ToMdmV1IDMConfig(i *v1.IDMConfig) *mdmv1.IDMConfig {
	if i == nil {
		return nil
	}
	return &mdmv1.IDMConfig{
		IdmType: i.IDMType,
	}
}

func ToV1Tenant(p *mdmv1.Tenant) *v1.Tenant {
	if p == nil {
		return nil
	}

	return &v1.Tenant{
		Meta:          ToV1Meta(p.Meta),
		Name:          p.Name,
		Description:   p.Description,
		Quotas:        ToV1QuotaSet(p.Quotas),
		DefaultQuotas: ToV1QuotaSet(p.DefaultQuotas),
		IAMConfig:     ToV1IAMConfig(p.IamConfig),
	}
}

func ToV1IAMConfig(i *mdmv1.IAMConfig) *v1.IAMConfig {
	if i == nil {
		return nil
	}

	return &v1.IAMConfig{
		IssuerConfig: ToV1IssuerConfig(i.IssuerConfig),
		IDMConfig:    ToV1IDMConfig(i.IdmConfig),
	}
}

func ToV1IssuerConfig(i *mdmv1.IssuerConfig) *v1.IssuerConfig {
	if i == nil {
		return nil
	}

	return &v1.IssuerConfig{
		URL:      i.Url,
		ClientID: i.ClientId,
	}
}

func ToV1IDMConfig(i *mdmv1.IDMConfig) *v1.IDMConfig {
	if i == nil {
		return nil
	}

	return &v1.IDMConfig{
		IDMType: i.IdmType,
	}
}

func ToMdmV1Project(p *v1.Project) *mdmv1.Project {
	if p == nil {
		return nil
	}

	return &mdmv1.Project{
		Meta:        ToMdmV1Meta(p.Meta),
		Name:        p.Name,
		Description: p.Description,
		TenantId:    p.TenantId,
		Quotas:      ToMdmV1QuotaSet(p.Quotas),
	}
}

func ToV1Project(p *mdmv1.Project) *v1.Project {
	if p == nil {
		return nil
	}

	return &v1.Project{
		Meta:        ToV1Meta(p.Meta),
		Name:        p.Name,
		Description: p.Description,
		TenantId:    p.TenantId,
		Quotas:      ToV1QuotaSet(p.Quotas),
	}
}

func ToMdmV1QuotaSet(qs *v1.QuotaSet) *mdmv1.QuotaSet {
	if qs == nil {
		return nil
	}

	return &mdmv1.QuotaSet{
		Cluster: ToMdmV1Quota(qs.Cluster),
		Machine: ToMdmV1Quota(qs.Machine),
		Ip:      ToMdmV1Quota(qs.Ip),
		Project: ToMdmV1Quota(qs.Project),
	}
}

func ToMdmV1Quota(q *v1.Quota) *mdmv1.Quota {
	if q == nil {
		return nil
	}
	if q.Quota == nil {
		return nil
	}

	return &mdmv1.Quota{
		Quota: &wrappers.Int32Value{
			Value: *q.Quota,
		},
	}
}

func ToMdmV1Meta(m *v1.Meta) *mdmv1.Meta {
	if m == nil {
		return nil
	}

	return &mdmv1.Meta{
		Id:          m.Id,
		Kind:        m.Kind,
		Apiversion:  m.Apiversion,
		Version:     m.Version,
		CreatedTime: mustTimeToTimestamp(m.CreatedTime),
		UpdatedTime: mustTimeToTimestamp(m.UpdatedTime),
		Annotations: m.Annotations,
		Labels:      m.Labels,
	}
}

func ToMdmV1ProjectFindRequest(r *v1.ProjectFindRequest) *mdmv1.ProjectFindRequest {
	if r == nil {
		return nil
	}

	mdmv1r := new(mdmv1.ProjectFindRequest)
	if r.Id != nil {
		mdmv1r.Id = &wrapperspb.StringValue{Value: *r.Id}
	}
	if r.Description != nil {
		mdmv1r.Description = &wrapperspb.StringValue{Value: *r.Description}
	}
	if r.Name != nil {
		mdmv1r.Name = &wrapperspb.StringValue{Value: *r.Name}
	}
	if r.TenantId != nil {
		mdmv1r.TenantId = &wrapperspb.StringValue{Value: *r.TenantId}
	}

	return mdmv1r
}

func ToV1Meta(m *mdmv1.Meta) *v1.Meta {
	if m == nil {
		return nil
	}

	return &v1.Meta{
		Id:          m.Id,
		Kind:        m.Kind,
		Apiversion:  m.Apiversion,
		Version:     m.Version,
		CreatedTime: mustTimestampToTime(m.CreatedTime),
		UpdatedTime: mustTimestampToTime(m.UpdatedTime),
		Annotations: m.Annotations,
		Labels:      m.Labels,
	}
}

func ToV1QuotaSet(q *mdmv1.QuotaSet) *v1.QuotaSet {
	if q == nil {
		return nil
	}
	return &v1.QuotaSet{
		Cluster: ToV1Quota(q.Cluster),
		Machine: ToV1Quota(q.Machine),
		Ip:      ToV1Quota(q.Ip),
		Project: ToV1Quota(q.Project),
	}
}

func ToV1Quota(q *mdmv1.Quota) *v1.Quota {
	if q == nil {
		return nil
	}
	return &v1.Quota{
		Quota: unwrapInt32(q.Quota),
	}
}
func unwrapInt32(w *wrappers.Int32Value) *int32 {
	if w == nil {
		return nil
	}

	return &w.Value
}

func mustTimestampToTime(ts *timestamp.Timestamp) *time.Time {
	t := ts.AsTime()
	return &t
}

func mustTimeToTimestamp(t *time.Time) *timestamp.Timestamp {
	return timestamppb.New(*t)
}
