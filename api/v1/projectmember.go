package v1

//go:generate go run ../../pkg/gen/genscanvaluer.go -package v1 -type ProjectMember

func (m *ProjectMember) NewProjectMemberResponse() *ProjectMemberResponse {
	return &ProjectMemberResponse{
		ProjectMember: m,
	}
}

func (m *ProjectMemberDeleteRequest) NewProjectMember() *ProjectMember {
	return &ProjectMember{
		Meta: &Meta{Id: m.Id},
	}
}
