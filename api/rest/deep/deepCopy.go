package deep

import (
	"github.com/metal-stack/masterdata-api/api/rest/mapper"
	mdmv1 "github.com/metal-stack/masterdata-api/api/v1"
)

func CopyTenant(t *mdmv1.Tenant) *mdmv1.Tenant {
	return mapper.ToMdmV1Tenant(mapper.ToV1Tenant(t))
}

func CopyProject(p *mdmv1.Project) *mdmv1.Project {
	return mapper.ToMdmV1Project(mapper.ToV1Project(p))
}
