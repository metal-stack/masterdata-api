package v1

import "time"

type (
	Tenant struct {
		Meta          *Meta      `json:"meta,omitempty"`
		Name          string     `json:"name,omitempty"`
		Description   string     `json:"description,omitempty"`
		DefaultQuotas *QuotaSet  `json:"default_quotas,omitempty"`
		Quotas        *QuotaSet  `json:"quotas,omitempty"`
		IAMConfig     *IAMConfig `json:"iam_config,omitempty"`
	}

	IAMConfig struct {
		IssuerConfig *IssuerConfig `json:"issuer_config,omitempty"`
		IDMConfig    *IDMConfig    `json:"idm_config,omitempty"`
	}

	IssuerConfig struct {
		URL      string `json:"url,omitempty"`
		ClientID string `json:"client_id,omitempty"`
	}

	IDMConfig struct {
		IDMType string `json:"idm_type,omitempty"`
	}

	TenantCreateRequest struct {
		Tenant
	}
	TenantUpdateRequest struct {
		Tenant
	}

	TenantFindRequest struct {
		Id          *string           `json:"id,omitempty"`
		Name        *string           `json:"name,omitempty"`
		Paging      *Paging           `json:"paging,omitempty"`
		Annotations map[string]string `json:"annotations,omitempty"`
	}

	TenantHistoryRequest struct {
		At time.Time `json:"at,omitempty"`
	}

	TenantResponse struct {
		Tenant
	}

	TenantListResponse struct {
		Tenants []*TenantResponse `json:"tenants,omitempty"`
	}

	Paging struct {
		Page  *uint64 `json:"page,omitempty"`
		Count *uint64 `json:"count,omitempty"`
	}
)
