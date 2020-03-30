package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/wrappers"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"testing"

	"github.com/metal-stack/masterdata-api/pkg/datastore/mocks"
)

func TestCreateProject(t *testing.T) {
	storageMock := &mocks.Storage{}
	ts := NewProjectService(storageMock, log)
	ctx := context.Background()

	t1 := &v1.Tenant{}
	p1 := &v1.Project{
		Name:        "FirstP",
		Description: "First Project",
		TenantId:    "t1",
		Meta: &v1.Meta{
			Annotations: map[string]string{
				"metal-stack.io/contract": "1234",
			},
			Labels: []string{
				"color=green",
			},
		},
	}
	tcr := &v1.ProjectCreateRequest{
		Project: p1,
	}
	storageMock.On("Get", ctx, p1.GetTenantId(), t1).Return(nil)
	storageMock.On("Create", ctx, p1).Return(nil)
	resp, err := ts.Create(ctx, tcr)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetProject())
	assert.Equal(t, tcr.Project.GetName(), resp.GetProject().GetName())
}

func TestCreateProjectWithQuotaCheck(t *testing.T) {
	storageMock := &mocks.Storage{}
	ts := NewProjectService(storageMock, log)
	ctx := context.Background()

	t1 := &v1.Tenant{
		Quotas: &v1.QuotaSet{
			Project: &v1.Quota{
				Quota: &wrappers.Int32Value{Value: 2},
			},
		},
	}
	p1 := &v1.Project{
		Name:        "FirstP",
		Description: "First Project",
		TenantId:    "t1",
	}
	tcr := &v1.ProjectCreateRequest{
		Project: p1,
	}
	filter := make(map[string]interface{})
	filter["project ->> 'tenant_id'"] = p1.TenantId
	var projects []v1.Project
	// see: https://github.com/stretchr/testify/blob/master/mock/mock.go#L149-L162
	storageMock.On("Get", ctx, p1.GetTenantId(), mock.AnythingOfType("*v1.Tenant")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*v1.Tenant)
		arg.Quotas = t1.GetQuotas()
	})
	storageMock.On("Find", ctx, filter, &projects).Return(nil)
	storageMock.On("Create", ctx, p1).Return(nil)
	resp, err := ts.Create(ctx, tcr)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetProject())
	assert.Equal(t, tcr.Project.GetName(), resp.GetProject().GetName())
}

func TestUpdateProject(t *testing.T) {
	storageMock := &mocks.Storage{}
	ts := NewProjectService(storageMock, log)
	ctx := context.Background()

	t1 := &v1.Project{
		Meta:        &v1.Meta{Id: "p2"},
		Name:        "SecondP",
		Description: "Second Project",
	}
	tur := &v1.ProjectUpdateRequest{
		Project: &v1.Project{
			Meta:        &v1.Meta{Id: "p2"},
			Name:        "SecondP",
			Description: "Second Project",
		},
	}

	storageMock.On("Update", ctx, t1).Return(nil)
	resp, err := ts.Update(ctx, tur)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetProject())
	assert.Equal(t, tur.GetProject().GetName(), resp.GetProject().GetName())
}

func TestDeleteProject(t *testing.T) {
	storageMock := &mocks.Storage{}
	ts := NewProjectService(storageMock, log)
	ctx := context.Background()
	t3 := &v1.Project{
		Meta: &v1.Meta{Id: "p3"},
	}
	tdr := &v1.ProjectDeleteRequest{
		Id: "p3",
	}

	storageMock.On("Delete", ctx, t3).Return(nil)
	resp, err := ts.Delete(ctx, tdr)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetProject())
	assert.Equal(t, tdr.Id, resp.GetProject().GetMeta().GetId())
}

// FIXME reenable
// func TestGetProject(t *testing.T) {
// 	storageMock := &mocks.Storage{}
// 	ts := NewProjectService(storageMock, log)
// 	ctx := context.Background()
// 	t4 := &v1.Project{}
// 	tgr := &v1.ProjectGetRequest{
// 		Id: "p4",
// 	}

// 	storageMock.On("Get", ctx, "p4", t4).Return(nil)
// 	resp, err := ts.Get(ctx, tgr)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	assert.NotNil(t, resp.GetProject())
// 	assert.Equal(t, tgr.Id, resp.GetProject().GetMeta().GetId())
// }

func TestFindProjectByID(t *testing.T) {
	storageMock := &mocks.Storage{}
	ts := NewProjectService(storageMock, log)
	ctx := context.Background()
	var t5s []v1.Project
	// filter by id
	f1 := make(map[string]interface{})
	tfr := &v1.ProjectFindRequest{
		Id: &wrappers.StringValue{Value: "p5"},
	}

	f1["id"] = "p5"
	storageMock.On("Find", ctx, f1, &t5s).Return(nil)
	resp, err := ts.Find(ctx, tfr)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestFindProjectByName(t *testing.T) {
	storageMock := &mocks.Storage{}
	ts := NewProjectService(storageMock, log)
	ctx := context.Background()

	// filter by name
	var t6s []v1.Project
	tfr := &v1.ProjectFindRequest{
		Name: &wrappers.StringValue{Value: "Sixth"},
	}

	f2 := make(map[string]interface{})
	f2["project ->> 'name'"] = "Sixth"
	storageMock.On("Find", ctx, f2, &t6s).Return(nil)
	resp, err := ts.Find(ctx, tfr)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestFindProjectByTenant(t *testing.T) {
	storageMock := &mocks.Storage{}
	ts := NewProjectService(storageMock, log)
	ctx := context.Background()

	// filter by name
	var t6s []v1.Project
	tfr := &v1.ProjectFindRequest{
		TenantId: &wrappers.StringValue{Value: "p1"},
	}

	f2 := make(map[string]interface{})
	f2["project ->> 'tenant_id'"] = "p1"
	storageMock.On("Find", ctx, f2, &t6s).Return(nil)
	resp, err := ts.Find(ctx, tfr)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
