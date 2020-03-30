package main

import (
	"context"
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

	c, err := client.NewClient(context.TODO(), "localhost", 50051, "certs/client.pem", "certs/client-key.pem", "certs/ca.pem", hmacKey, logger)
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
		log.Sugar().Fatalf("could not create project: %v", err)
	}
	log.Sugar().Infow("created project", "project", res)

	// get
	prj, err := c.Project().Get(ctx, &v1.ProjectGetRequest{Id: res.Project.Meta.Id})
	if err != nil {
		log.Sugar().Fatalf("created project notfound, id=%s", res.Project.Meta.Id)
	}

	// update
	prj.Project.Meta.Annotations["mykey"] = "myvalue"
	prures, err := c.Project().Update(ctx, &v1.ProjectUpdateRequest{
		Project: prj.Project,
	})
	if err != nil {
		log.Sugar().Fatalf("update project failed, id=%s", res.Project.Meta.Id)
	}

	// explicit re-get
	prj2, err := c.Project().Get(ctx, &v1.ProjectGetRequest{Id: res.Project.Meta.Id})
	if err != nil {
		log.Sugar().Fatalf("created project notfound, id=%s", res.Project.Meta.Id)
	}
	if prj2.GetProject().Meta.Annotations["mykey"] != "myvalue" {
		log.Sugar().Fatalf("update project failed, id=%s", res.Project.Meta.Id)
	}

	if prures.Project.Meta.Version <= prj.Project.Meta.Version {
		log.Sugar().Fatalf("update project failed, version not incremented", res.Project.Meta.Id)
	}

	_, err = c.Project().Get(ctx, &v1.ProjectGetRequest{Id: "123123"})
	if !v1.IsNotFound(err) {
		log.Sugar().Fatalf("expected notfound")
	}

	// find
	pfr, err := c.Project().Find(ctx, &v1.ProjectFindRequest{})
	if err != nil {
		log.Sugar().Fatalf("could get create find projects endpoint: %v", err)
	}
	for _, p := range pfr.Projects {
		log.Sugar().Infow("found project", "project", p)
	}

	// delete projects
	for _, p := range pfr.GetProjects() {
		pdr := v1.ProjectDeleteRequest{
			Id: p.Meta.Id,
		}
		_, err = c.Project().Delete(ctx, &pdr)
		if err != nil {
			log.Sugar().Fatalf("could delete project: %v", err)
		}
		log.Sugar().Infow("deleted ", "project", p)
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
				Url:      "https://dex.fi-ts.io/dex",
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
		log.Sugar().Fatalf("could not create tenant: %v", err)
	}
	log.Sugar().Infof("created tenant: %v", t)

	// try to create the same tenant with the returned id another time...
	tcr.Tenant.Meta = t.Tenant.Meta
	_, err = c.Tenant().Create(ctx, tcr)

	if err != nil {
		if v1.IsConflict(err) {
			log.Sugar().Info("got expected grpc code, indicating duplicate key")
		} else {
			log.Sugar().Fatalf("could not create tenant, unexpected error: %v", err)
		}
	} else {
		log.Sugar().Fatalf("THIS MUST NOT HAPPEN: successfully created tenant with duplicate ID: %v", t)
	}

	log.Sugar().Infof("find all tenants")
	tfrq := &v1.TenantFindRequest{
		// Id:                   t.Id,
	}
	tfrs, err := c.Tenant().Find(ctx, tfrq)
	if err != nil {
		log.Sugar().Fatalf("could not find tenants: %v", err)
	}
	for i := range tfrs.Tenants {
		log.Sugar().Infof("found tenant %v", tfrs.Tenants[i])
	}

	log.Sugar().Infof("get tenant with id")
	tgr := &v1.TenantGetRequest{
		Id: t.Tenant.Meta.Id,
	}
	tgres, err := c.Tenant().Get(ctx, tgr)
	if err != nil {
		log.Sugar().Fatalf("could not get tenant: %v", err)
	}
	log.Sugar().Infof("got tenant with id %v", tgres)

	log.Sugar().Infof("get tenant with non-existant id")
	tgrNotFound := &v1.TenantGetRequest{
		Id: "1982739817298219873",
	}
	_, err = c.Tenant().Get(ctx, tgrNotFound)
	if !v1.IsNotFound(err) {
		log.Sugar().Fatal("unexpected response %v with error on tenant that cannot be found!", err)
	}

	// get tenant one more time to have some older version after update to provoke an optimistic lock error
	tgres2, _ := c.Tenant().Get(ctx, tgr)

	tenant := tgres.Tenant
	tenant.Name = "some other name"

	tur := &v1.TenantUpdateRequest{
		Tenant: tenant,
	}
	tures, err := c.Tenant().Update(ctx, tur)
	if err != nil {
		log.Sugar().Fatalf("could not update tenant: %v", err)
	}
	log.Sugar().Infof("updated tenant %v", tures)

	tenant2 := tgres2.Tenant
	tenant2.Name = "update older tenant"
	_, err = c.Tenant().Update(ctx, tur)
	if !v1.IsOptimistickLockError(err) {
		log.Sugar().Fatalf("could not update tenant, expected OptimistickLockError, got error: %v", err)
	}

	log.Sugar().Infof("find tenant with id")
	tfrqi := &v1.TenantFindRequest{
		Id: &wrappers.StringValue{Value: t.Tenant.Meta.Id},
	}
	tfrsi, err := c.Tenant().Find(ctx, tfrqi)
	if err != nil {
		log.Sugar().Fatalf("could not find tenants: %v", err)
	}
	for i := range tfrsi.Tenants {
		log.Sugar().Infof("found tenant %v", tfrsi.Tenants[i])
	}

	log.Sugar().Infof("delete tenant with id")
	tdr := &v1.TenantDeleteRequest{
		Id: t.Tenant.Meta.Id,
	}
	_, err = c.Tenant().Delete(ctx, tdr)
	if err != nil {
		log.Sugar().Fatalf("could not delete tenant: %v", err)
	}

	log.Sugar().Infof("try to delete already deleted tenant")
	tdr2 := &v1.TenantDeleteRequest{
		Id: t.Tenant.Meta.Id,
	}
	_, err = c.Tenant().Delete(ctx, tdr2)
	if !v1.IsNotFound(err) {
		log.Sugar().Info("got expected grpc code, indicating not found")
	}
}
