package v1

//go:generate go run ../../pkg/gen/genscanvaluer.go -package v1 -type TenantMember

func (m *TenantMember) NewTenantMemberResponse() *TenantMemberResponse {
	return &TenantMemberResponse{
		TenantMember: m,
	}
}

func (m *TenantMemberDeleteRequest) NewTenantMember() *TenantMember {
	return &TenantMember{
		Meta: &Meta{Id: m.Id},
	}
}
