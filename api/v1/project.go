package v1

//go:generate go run ../../pkg/gen/genscanvaluer.go -package v1 -type Project

func (m *Project) NewProjectResponse() *ProjectResponse {
	return &ProjectResponse{
		Project: m,
	}
}

func (m *ProjectUpdateRequest) NewProject() *Project {
	return &Project{
		Meta:        &Meta{Id: m.Id},
		Name:        m.GetName().GetValue(),
		Description: m.GetDescription().GetValue(),
		Quotas:      m.Quotas,
	}
}

func (m *ProjectDeleteRequest) NewProject() *Project {
	return &Project{
		Meta: &Meta{Id: m.Id},
	}
}
