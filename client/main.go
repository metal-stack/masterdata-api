package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/metal-stack/masterdata-api/api/rest/mapper"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/metal-stack/masterdata-api/pkg/auth"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/client"
)

const grpcRequestTimeout = 5 * time.Second

func main() {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})

	logger := slog.New(jsonHandler)
	logger.Info("Starting Client")

	hmacKey := os.Getenv("HMAC_KEY")
	if hmacKey == "" {
		hmacKey = auth.HmacDefaultKey
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, err := client.NewClient(ctx, "localhost", 50051, "certs/client.pem", "certs/client-key.pem", "certs/ca.pem", hmacKey, true, logger)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	defer func() {
		err = c.Close()
		if err != nil {
			logger.Error(err.Error())
			panic(err)
		}
	}()
	err = projectExample(c, logger)
	if err != nil {
		logger.Error(err.Error())
	}
	err = tenantExample(c, logger)
	if err != nil {
		logger.Error(err.Error())
	}

	logger.Info("Success")
}

func projectExample(c client.Client, log *slog.Logger) error {

	ctx, cancel := context.WithTimeout(context.Background(), grpcRequestTimeout)
	defer cancel()

	// create
	project := &v1.Project{
		Name:        "project123",
		Description: "Demo Project",
		TenantId:    "customer-1",
		Quotas: &v1.QuotaSet{
			Cluster: &v1.Quota{Quota: pointer.Pointer(int32(3))},
			Machine: &v1.Quota{Quota: pointer.Pointer(int32(3))},
			Ip:      &v1.Quota{Quota: pointer.Pointer(int32(3))},
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
		return err
	}
	log.Info("created project", "project", res)

	ts := time.Now()

	// get
	projectId := res.Project.Meta.Id
	prj, err := c.Project().Get(ctx, &v1.ProjectGetRequest{Id: projectId})
	if err != nil {
		return err
	}

	// update
	prj.Project.Description = "Updated Demo Project"
	prj.Project.Meta.Annotations["mykey"] = "myvalue"
	prures, err := c.Project().Update(ctx, &v1.ProjectUpdateRequest{
		Project: prj.Project,
	})
	if err != nil {
		return err
	}

	pbHp := timestamppb.New(ts)
	phr, err := c.Project().GetHistory(ctx, &v1.ProjectGetHistoryRequest{
		Id: prj.Project.Meta.Id,
		At: pbHp,
	})
	if err != nil {
		return err
	}
	if phr.Project.Description != "Demo Project" {
		return err
	}

	// explicit re-get
	prj2, err := c.Project().Get(ctx, &v1.ProjectGetRequest{Id: projectId})
	if err != nil {
		return err
	}
	if prj2.GetProject().Meta.Annotations["mykey"] != "myvalue" {
		return err
	}

	if prures.Project.Meta.Version <= prj.Project.Meta.Version {
		return err
	}

	_, err = c.Project().Get(ctx, &v1.ProjectGetRequest{Id: "123123"})
	if !v1.IsNotFound(err) {
		return err
	}

	// find
	pfr, err := c.Project().Find(ctx, &v1.ProjectFindRequest{})
	if err != nil {
		return err
	}
	for _, p := range pfr.Projects {
		log.Info("found project", "project", p)
	}

	pmcr, err := c.ProjectMember().Create(ctx, &v1.ProjectMemberCreateRequest{
		ProjectMember: &v1.ProjectMember{
			ProjectId: projectId,
			TenantId:  "customer-1",
		},
	})
	if err != nil {
		return err
	}

	log.Info("projectmember created", slog.Any("member", pmcr.ProjectMember))

	// delete projects
	for _, p := range pfr.GetProjects() {
		pdr := v1.ProjectDeleteRequest{
			Id: p.Meta.Id,
		}
		_, err = c.Project().Delete(ctx, &pdr)
		if err != nil {
			return err
		}
		log.Info("deleted ", "project", p)
	}
	return nil
}

func tenantExample(c client.Client, log *slog.Logger) error {
	tnt := &v1.Tenant{
		Meta:        nil,
		Name:        "myTenant",
		Description: "myDesc",
		DefaultQuotas: &v1.QuotaSet{
			Cluster: &v1.Quota{Quota: pointer.Pointer(int32(3))},
			Machine: &v1.Quota{Quota: pointer.Pointer(int32(3))},
			Ip:      &v1.Quota{Quota: pointer.Pointer(int32(3))},
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
		return err
	}
	log.Info("created tenant", "tenant", t)

	// try to create the same tenant with the returned id another time...
	tcr.Tenant.Meta = t.Tenant.Meta
	_, err = c.Tenant().Create(ctx, tcr)

	if err != nil {
		if v1.IsConflict(err) {
			log.Info("got expected grpc code, indicating duplicate key")
		} else {
			return err
		}
	} else {
		return err
	}

	log.Info("find all tenants")
	tfrq := &v1.TenantFindRequest{
		// Id:                   t.Id,
	}
	tfrs, err := c.Tenant().Find(ctx, tfrq)
	if err != nil {
		return err
	}
	for i := range tfrs.Tenants {
		log.Info("found tenant", "tenant", tfrs.Tenants[i])
	}

	log.Info("get tenant with id")
	tgr := &v1.TenantGetRequest{
		Id: t.Tenant.Meta.Id,
	}
	tgres, err := c.Tenant().Get(ctx, tgr)
	if err != nil {
		return err
	}
	log.Info("got tenant", "id", tgres)

	v1t := mapper.ToV1Tenant(tgres.Tenant)
	mdm1t := mapper.ToMdmV1Tenant(v1t)

	_, err = c.Tenant().Update(ctx, &v1.TenantUpdateRequest{Tenant: mdm1t})
	if err != nil {
		return err

	}

	log.Info("get tenant with non-existent id")
	tgrNotFound := &v1.TenantGetRequest{
		Id: "1982739817298219873",
	}
	_, err = c.Tenant().Get(ctx, tgrNotFound)
	if !v1.IsNotFound(err) {
		return err
	}

	// get tenant one more time to have some older version after update to provoke an optimistic lock error
	tgres2, _ := c.Tenant().Get(ctx, tgr)

	tgres, err = c.Tenant().Get(ctx, tgr)
	if err != nil {
		return err
	}
	tenant := tgres.Tenant
	tenant.Name = "some other name"

	tur := &v1.TenantUpdateRequest{
		Tenant: tenant,
	}
	tures, err := c.Tenant().Update(ctx, tur)
	if err != nil {
		return err
	}
	log.Info("updated tenant", "tenant", tures)

	tenant2 := tgres2.Tenant
	tenant2.Name = "update older tenant"
	_, err = c.Tenant().Update(ctx, tur)
	if !v1.IsOptimistickLockError(err) {
		return err
	}

	log.Info("find tenant with id")
	tfrqi := &v1.TenantFindRequest{
		Id: pointer.Pointer(t.Tenant.Meta.Id),
	}
	tfrsi, err := c.Tenant().Find(ctx, tfrqi)
	if err != nil {
		return err
	}
	for i := range tfrsi.Tenants {
		log.Info("found tenant", "tenant", tfrsi.Tenants[i])
	}

	log.Info("delete tenant with id")
	tdr := &v1.TenantDeleteRequest{
		Id: t.Tenant.Meta.Id,
	}
	_, err = c.Tenant().Delete(ctx, tdr)
	if err != nil {
		return err
	}

	log.Info("try to delete already deleted tenant")
	tdr2 := &v1.TenantDeleteRequest{
		Id: t.Tenant.Meta.Id,
	}
	_, err = c.Tenant().Delete(ctx, tdr2)
	if !v1.IsNotFound(err) {
		log.Info("got expected grpc code, indicating not found")
	}

	pbHt := timestamppb.Now()
	thr, err := c.Tenant().GetHistory(ctx, &v1.TenantGetHistoryRequest{
		Id: tdr.Id,
		At: pbHt,
	})
	if err != nil {
		return err
	}
	if thr.Tenant.Name != "some other name" {
		return err
	}
	log.Info("found history tenant", "tenant", thr.Tenant)
	return nil
}
