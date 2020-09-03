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

	ProjectFindRequest struct {
		Id          *string `json:"id,omitempty"`
		Name        *string `json:"name,omitempty"`
		Description *string `json:"description,omitempty"`
		TenantId    *string `json:"tenant_id,omitempty"`
	}

	ProjectResponse struct {
		Project
	}

	ProjectListResponse struct {
		Projects []*ProjectResponse `json:"projects,omitempty"`
	}
)
