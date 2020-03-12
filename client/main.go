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

const grpcRequestTimeoutSecs = 5 * time.Second

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

	ctx, cancel := context.WithTimeout(context.Background(), grpcRequestTimeoutSecs)
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
	}
	pcr := &v1.ProjectCreateRequest{
		Project: project,
	}
	res, err := c.Project().Create(ctx, pcr)
	if err != nil {
		log.Sugar().Fatal("could not create project", zap.Error(err))
	}
	log.Sugar().Infow("created project", "project", res)

	// get
	_, err = c.Project().Get(ctx, &v1.ProjectGetRequest{Id: res.Project.Meta.Id})
	if err != nil {
		log.Sugar().Fatal("created project notfound", "id", res.Project.Meta.Id)
	}
	_, err = c.Project().Get(ctx, &v1.ProjectGetRequest{Id: "123123"})
	if !v1.IsNotFound(err) {
		log.Sugar().Fatal("expected notfound")
	}

	// find
	pfr, err := c.Project().Find(ctx, &v1.ProjectFindRequest{})
	if err != nil {
		log.Sugar().Fatal("could get create find projects endpoint", zap.Error(err))
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
			log.Sugar().Fatal("could delete project", zap.Error(err))
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
	ctx, cancel := context.WithTimeout(context.Background(), grpcRequestTimeoutSecs*time.Second)
	defer cancel()

	tcr := &v1.TenantCreateRequest{
		Tenant: tnt,
	}

	t, err := c.Tenant().Create(ctx, tcr)
	if err != nil {
		log.Sugar().Fatal("could not create tenant", zap.Error(err))
	}
	log.Sugar().Info("created tenant", "tenant", t)

	// try to create the same tenant with the returned id another time...
	tcr.Tenant.Meta = t.Tenant.Meta
	_, err = c.Tenant().Create(ctx, tcr)

	if err != nil {
		if v1.IsConflict(err) {
			log.Sugar().Info("got expected grpc code, indicating duplicate key")
		} else {
			log.Sugar().Fatal("could not create tenant, unexpected error", zap.Error(err))
		}
	} else {
		log.Sugar().Fatal("THIS MUST NOT HAPPEN: successfully created tenant with duplicate ID", "ID", t)
	}

	log.Sugar().Info("find all tenants")
	tfrq := &v1.TenantFindRequest{
		// Id:                   t.Id,
	}
	tfrs, err := c.Tenant().Find(ctx, tfrq)
	if err != nil {
		log.Sugar().Fatal("could not find tenants", zap.Error(err))
	}
	for i := range tfrs.Tenants {
		log.Sugar().Info("found tenant", "tenant", tfrs.Tenants[i])
	}

	log.Sugar().Info("get tenant with id")
	tgr := &v1.TenantGetRequest{
		Id: t.Tenant.Meta.Id,
	}
	tgres, err := c.Tenant().Get(ctx, tgr)
	if err != nil {
		log.Sugar().Fatal("could not get tenant", zap.Error(err))
	}
	log.Sugar().Info("got tenant", "id", tgres)

	log.Sugar().Info("get tenant with non-existant id")
	tgr_notfound := &v1.TenantGetRequest{
		Id: "1982739817298219873",
	}
	_, err = c.Tenant().Get(ctx, tgr_notfound)
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
		log.Sugar().Fatal("could not update tenant", zap.Error(err))
	}
	log.Sugar().Info("updated tenant", "tenant", tures)

	tenant2 := tgres2.Tenant
	tenant2.Name = "update older tenant"
	_, err = c.Tenant().Update(ctx, tur)
	if !v1.IsOptimistickLockError(err) {
		log.Sugar().Fatal("could not update tenant, expected OptimistickLockError, got error", zap.Error(err))
	}

	log.Sugar().Info("find tenant with id")
	tfrqi := &v1.TenantFindRequest{
		Id: &wrappers.StringValue{Value: t.Tenant.Meta.Id},
	}
	tfrsi, err := c.Tenant().Find(ctx, tfrqi)
	if err != nil {
		log.Sugar().Fatal("could not find tenants", zap.Error(err))
	}
	for i := range tfrsi.Tenants {
		log.Sugar().Info("found tenant", "tenant", tfrsi.Tenants[i])
	}

	log.Sugar().Info("delete tenant with id")
	tdr := &v1.TenantDeleteRequest{
		Id: t.Tenant.Meta.Id,
	}
	_, err = c.Tenant().Delete(ctx, tdr)
	if err != nil {
		log.Sugar().Fatal("could not delete tenant", zap.Error(err))
	}

	log.Sugar().Info("try to delete already deleted tenant")
	tdr2 := &v1.TenantDeleteRequest{
		Id: t.Tenant.Meta.Id,
	}
	_, err = c.Tenant().Delete(ctx, tdr2)
	if !v1.IsNotFound(err) {
		log.Sugar().Info("got expected grpc code, indicating not found")
	}
}
