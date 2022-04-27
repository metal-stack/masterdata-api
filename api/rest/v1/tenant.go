package v1

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

	TenantUpdateRequest struct {
		Tenant *Tenant `json:"tenant,omitempty"`
	}

	TenantFindRequest struct {
		Id   *string `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
	}

	TenantResponse struct {
		Tenant
	}

	TenantListResponse struct {
		Tenants []*TenantResponse `json:"tenants,omitempty"`
	}
)
