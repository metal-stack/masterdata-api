package v1

//go:generate go run ../../pkg/gen/genscanvaluer.go -package v1 -type Tenant

func (m *Tenant) NewTenantResponse() *TenantResponse {
	return &TenantResponse{
		Tenant: m,
	}
}

func (m *TenantDeleteRequest) NewTenant() *Tenant {
	return &Tenant{
		Meta: &Meta{Id: m.Id},
	}
}
