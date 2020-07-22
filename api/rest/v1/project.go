package v1

type (
	Project struct {
		Meta        *Meta     `json:"meta,omitempty"`
		Name        string    `json:"name,omitempty"`
		Description string    `json:"description,omitempty"`
		TenantId    string    `json:"tenant_id,omitempty"`
		Quotas      *QuotaSet `json:"quotas,omitempty"`
	}

	ProjectCreateRequest struct {
		Project
	}

	ProjectUpdateRequest struct {
		Project
	}

	ProjectResponse struct {
		Project
	}

	ProjectListResponse struct {
		Projects []*Project `json:"projects,omitempty"`
	}
)
