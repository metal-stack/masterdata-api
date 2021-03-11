package main

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/metal-stack/masterdata-api/api/rest/mapper"
	"os"
	"time"

	"github.com/metal-stack/masterdata-api/pkg/auth"

	"github.com/golang/protobuf/ptypes/wrappers"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/client"
	"go.uber.org/zap"
)

const grpcRequestTimeout = 5 * time.Second

func main() {

	logger, _ := zap.NewProduction()
	logger.Info("Starting Client")

	hmacKey := os.Getenv("HMAC_KEY")
	if hmacKey == "" {
		hmacKey = auth.HmacDefaultKey
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	c, err := client.NewClient(ctx, "localhost", 50051, "certs/client.pem", "certs/client-key.pem", "certs/ca.pem", hmacKey, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer func() {
		err = c.Close()
		if err != nil {
			logger.Fatal(err.Error())
		}
	}()
	projectExample(c, logger)
	tenantExample(c, logger)

	logger.Info("Success")
}

func projectExample(c client.Client, log *zap.Logger) {

	ctx, cancel := context.WithTimeout(context.Background(), grpcRequestTimeout)
	defer cancel()

	// create
	project := &v1.Project{
		Name:        "project123",
		Description: "Demo Project",
		TenantId:    "customer-1",
		Quotas: &v1.QuotaSet{
			Cluster: &v1.Quota{Quota: &wrappers.Int32Value{Value: 3}},
			Machine: &v1.Quota{Quota: &wrappers.Int32Value{Value: 3}},
			Ip:      &v1.Quota{Quota: &wrappers.Int32Value{Value: 3}},
		},
		Meta: &v1.Meta{
			Annotations: map[string]string{
				"metal-stack.io/contract": "1234",
			},
			Labels: []string{
				"color=green",
			},
		},
	}
	pcr := &v1.ProjectCreateRequest{
		Project: project,
	}
	res, err := c.Project().Create(ctx, pcr)
	if err != nil {
		log.Fatal("could not create project", zap.Error(err))
	}
	log.Info("created project", zap.Stringer("project", res))

	ts := time.Now()

	// get
	projectId := res.Project.Meta.Id
	prj, err := c.Project().Get(ctx, &v1.ProjectGetRequest{Id: projectId})
	if err != nil {
		log.Fatal("created project notfound", zap.String("id", projectId))
	}

	// update
	prj.Project.Description = "Updated Demo Project"
	prj.Project.Meta.Annotations["mykey"] = "myvalue"
	prures, err := c.Project().Update(ctx, &v1.ProjectUpdateRequest{
		Project: prj.Project,
	})
	if err != nil {
		log.Fatal("update project failed", zap.String("id", projectId))
	}

	pbHp, _ := ptypes.TimestampProto(ts)
	phr, err := c.Project().GetHistory(ctx, &v1.ProjectGetHistoryRequest{
		Id: prj.Project.Meta.Id,
		At: pbHp,
	})
	if err != nil {
		log.Fatal("get project history failed", zap.String("id", projectId))
	}
	if phr.Project.Description != "Demo Project" {
		log.Fatal("get project: unexpected description", zap.String("id", projectId), zap.String("desc", phr.Project.Description))
	}

	// explicit re-get
	prj2, err := c.Project().Get(ctx, &v1.ProjectGetRequest{Id: projectId})
	if err != nil {
		log.Fatal("created project notfound", zap.String("id", projectId))
	}
	if prj2.GetProject().Meta.Annotations["mykey"] != "myvalue" {
		log.Fatal("update project failed", zap.String("id", projectId))
	}

	if prures.Project.Meta.Version <= prj.Project.Meta.Version {
		log.Fatal("update project failed, version not incremented", zap.String("id", projectId))
	}

	_, err = c.Project().Get(ctx, &v1.ProjectGetRequest{Id: "123123"})
	if !v1.IsNotFound(err) {
		log.Fatal("expected notfound")
	}

	// find
	pfr, err := c.Project().Find(ctx, &v1.ProjectFindRequest{})
	if err != nil {
		log.Fatal("could get create find projects endpoint", zap.Error(err))
	}
	for _, p := range pfr.Projects {
		log.Info("found project", zap.Stringer("project", p))
	}

	// delete projects
	for _, p := range pfr.GetProjects() {
		pdr := v1.ProjectDeleteRequest{
			Id: p.Meta.Id,
		}
		_, err = c.Project().Delete(ctx, &pdr)
		if err != nil {
			log.Fatal("could delete project", zap.Error(err))
		}
		log.Info("deleted ", zap.Stringer("project", p))
	}
}

func tenantExample(c client.Client, log *zap.Logger) {
	tnt := &v1.Tenant{
		Meta:        nil,
		Name:        "myTenant",
		Description: "myDesc",
		DefaultQuotas: &v1.QuotaSet{
			Cluster: &v1.Quota{Quota: &wrappers.Int32Value{Value: 3}},
			Machine: &v1.Quota{Quota: &wrappers.Int32Value{Value: 3}},
			Ip:      &v1.Quota{Quota: &wrappers.Int32Value{Value: 3}},
		},
		IamConfig: &v1.IAMConfig{
			IssuerConfig: &v1.IssuerConfig{
				Url:      "https://dex.test.metal-stack.io/dex",
				ClientId: "123213213",
			},
			IdmConfig: &v1.IDMConfig{
				IdmType: "UX",
				ConnectorConfig: &v1.ConnectorConfig{
					IdmApiUrl:      "a",
					IdmApiUser:     "b",
					IdmApiPassword: "c",
					IdmSystemId:    "d",
					IdmAccessCode:  "e",
					IdmCustomerId:  "f",
					IdmGroupOu:     "g",
					IdmGroupnameTemplate: &wrappers.StringValue{
						Value: "asdasdads",
					},
				},
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), grpcRequestTimeout)
	defer cancel()

	tcr := &v1.TenantCreateRequest{
		Tenant: tnt,
	}

	t, err := c.Tenant().Create(ctx, tcr)
	if err != nil {
		log.Fatal("could not create tenant", zap.Error(err))
	}
	log.Info("created tenant", zap.Stringer("tenant", t))

	// try to create the same tenant with the returned id another time...
	tcr.Tenant.Meta = t.Tenant.Meta
	_, err = c.Tenant().Create(ctx, tcr)

	if err != nil {
		if v1.IsConflict(err) {
			log.Info("got expected grpc code, indicating duplicate key")
		} else {
			log.Fatal("could not create tenant, unexpected error", zap.Error(err))
		}
	} else {
		log.Fatal("THIS MUST NOT HAPPEN: successfully created tenant with duplicate ID", zap.Stringer("ID", t))
	}

	log.Info("find all tenants")
	tfrq := &v1.TenantFindRequest{
		// Id:                   t.Id,
	}
	tfrs, err := c.Tenant().Find(ctx, tfrq)
	if err != nil {
		log.Fatal("could not find tenants", zap.Error(err))
	}
	for i := range tfrs.Tenants {
		log.Info("found tenant", zap.Stringer("tenant", tfrs.Tenants[i]))
	}

	log.Info("get tenant with id")
	tgr := &v1.TenantGetRequest{
		Id: t.Tenant.Meta.Id,
	}
	tgres, err := c.Tenant().Get(ctx, tgr)
	if err != nil {
		log.Fatal("could not get tenant", zap.Error(err))
	}
	log.Info("got tenant", zap.Stringer("id", tgres))

	v1t := mapper.ToV1Tenant(tgres.Tenant)
	mdm1t := mapper.ToMdmV1Tenant(v1t)

	_, err = c.Tenant().Update(ctx, &v1.TenantUpdateRequest{Tenant: mdm1t})
	if err != nil {
		log.Fatal("could not get tenant", zap.Error(err))
	}

	log.Info("get tenant with non-existant id")
	tgrNotFound := &v1.TenantGetRequest{
		Id: "1982739817298219873",
	}
	_, err = c.Tenant().Get(ctx, tgrNotFound)
	if !v1.IsNotFound(err) {
		log.Fatal("unexpected response with error on tenant that cannot be found!", zap.Error(err))
	}

	// get tenant one more time to have some older version after update to provoke an optimistic lock error
	tgres2, _ := c.Tenant().Get(ctx, tgr)

	tgres, err = c.Tenant().Get(ctx, tgr)
	if err != nil {
		log.Fatal("could not get tenant", zap.Error(err))
	}
	tenant := tgres.Tenant
	tenant.Name = "some other name"

	tur := &v1.TenantUpdateRequest{
		Tenant: tenant,
	}
	tures, err := c.Tenant().Update(ctx, tur)
	if err != nil {
		log.Fatal("could not update tenant", zap.Error(err))
	}
	log.Info("updated tenant", zap.Stringer("tenant", tures))

	tenant2 := tgres2.Tenant
	tenant2.Name = "update older tenant"
	_, err = c.Tenant().Update(ctx, tur)
	if !v1.IsOptimistickLockError(err) {
		log.Fatal("could not update tenant, expected OptimisticLockError, got error", zap.Error(err))
	}

	log.Info("find tenant with id")
	tfrqi := &v1.TenantFindRequest{
		Id: &wrappers.StringValue{Value: t.Tenant.Meta.Id},
	}
	tfrsi, err := c.Tenant().Find(ctx, tfrqi)
	if err != nil {
		log.Fatal("could not find tenants", zap.Error(err))
	}
	for i := range tfrsi.Tenants {
		log.Info("found tenant", zap.Stringer("tenant", tfrsi.Tenants[i]))
	}

	log.Info("delete tenant with id")
	tdr := &v1.TenantDeleteRequest{
		Id: t.Tenant.Meta.Id,
	}
	_, err = c.Tenant().Delete(ctx, tdr)
	if err != nil {
		log.Fatal("could not delete tenant", zap.Error(err))
	}

	log.Info("try to delete already deleted tenant")
	tdr2 := &v1.TenantDeleteRequest{
		Id: t.Tenant.Meta.Id,
	}
	_, err = c.Tenant().Delete(ctx, tdr2)
	if !v1.IsNotFound(err) {
		log.Info("got expected grpc code, indicating not found")
	}

	pbHt, _ := ptypes.TimestampProto(time.Now())
	thr, err := c.Tenant().GetHistory(ctx, &v1.TenantGetHistoryRequest{
		Id: tdr.Id,
		At: pbHt,
	})
	if err != nil {
		log.Fatal("tenant history not found", zap.Error(err))
	}
	if thr.Tenant.Name != "some other name" {
		log.Fatal("get tenant: unexpected name", zap.String("id", tdr.Id), zap.String("name", thr.Tenant.Name))
	}
	log.Info("found history tenant", zap.Stringer("tenant", thr.Tenant))
}
