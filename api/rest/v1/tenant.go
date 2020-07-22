package v1

type (
	Tenant struct {
		Meta          *Meta     `json:"meta,omitempty"`
		Name          string    `json:"name,omitempty"`
		Description   string    `json:"description,omitempty"`
		DefaultQuotas *QuotaSet `json:"default_quotas,omitempty"`
		Quotas        *QuotaSet `json:"quotas,omitempty"`
	}

	TenantUpdateRequest struct {
		Tenant *Tenant `json:"tenant,omitempty"`
	}

	TenantResponse struct {
		Tenant *Tenant `json:"tenant,omitempty"`
	}
)
