package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"connectrpc.com/connect"
	"github.com/metal-stack/masterdata-api/api/rest/mapper"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/client"
)

const grpcRequestTimeout = 5 * time.Second

func main() {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})

	logger := slog.New(jsonHandler)
	logger.Info("Starting Client")

	c := client.New(client.DialConfig{
		BaseURL:   "http://localhost:9090",
		Debug:     true,
		UserAgent: "sample-client",
	})

	err := projectExample(c, logger)
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
	res, err := c.Project().Create(ctx, connect.NewRequest(pcr))
	if err != nil {
		return err
	}
	log.Info("created project", "project", res)

	ts := time.Now()

	// get
	projectId := res.Msg.Project.Meta.Id
	prj, err := c.Project().Get(ctx, connect.NewRequest(&v1.ProjectGetRequest{Id: projectId}))
	if err != nil {
		return err
	}

	// update
	prj.Msg.Project.Description = "Updated Demo Project"
	prj.Msg.Project.Meta.Annotations["mykey"] = "myvalue"
	prures, err := c.Project().Update(ctx, connect.NewRequest(&v1.ProjectUpdateRequest{
		Project: prj.Msg.Project,
	}))
	if err != nil {
		return err
	}

	pbHp := timestamppb.New(ts)
	phr, err := c.Project().GetHistory(ctx, connect.NewRequest(&v1.ProjectGetHistoryRequest{
		Id: prj.Msg.Project.Meta.Id,
		At: pbHp,
	}))
	if err != nil {
		return err
	}
	if phr.Msg.Project.Description != "Demo Project" {
		return err
	}

	// explicit re-get
	prj2, err := c.Project().Get(ctx, connect.NewRequest(&v1.ProjectGetRequest{Id: projectId}))
	if err != nil {
		return err
	}
	if prj2.Msg.GetProject().Meta.Annotations["mykey"] != "myvalue" {
		return err
	}

	if prures.Msg.Project.Meta.Version <= prj.Msg.Project.Meta.Version {
		return err
	}

	_, err = c.Project().Get(ctx, connect.NewRequest(&v1.ProjectGetRequest{Id: "123123"}))
	if !v1.IsNotFound(err) {
		return err
	}

	// find
	pfr, err := c.Project().Find(ctx, connect.NewRequest(&v1.ProjectFindRequest{}))
	if err != nil {
		return err
	}
	for _, p := range pfr.Msg.Projects {
		log.Info("found project", "project", p)
	}

	pmcr, err := c.ProjectMember().Create(ctx, connect.NewRequest(&v1.ProjectMemberCreateRequest{
		ProjectMember: &v1.ProjectMember{
			ProjectId: projectId,
			TenantId:  "customer-1",
		},
	}))
	if err != nil {
		return err
	}

	log.Info("projectmember created", slog.Any("member", pmcr.Msg.ProjectMember))

	// delete projects
	for _, p := range pfr.Msg.GetProjects() {
		pdr := v1.ProjectDeleteRequest{
			Id: p.Meta.Id,
		}
		_, err = c.Project().Delete(ctx, connect.NewRequest(&pdr))
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
					IdmApiUrl:            "a",
					IdmApiUser:           "b",
					IdmApiPassword:       "c",
					IdmSystemId:          "d",
					IdmAccessCode:        "e",
					IdmCustomerId:        "f",
					IdmGroupOu:           "g",
					IdmGroupnameTemplate: pointer.Pointer("asdasdads"),
				},
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), grpcRequestTimeout)
	defer cancel()

	tcr := &v1.TenantCreateRequest{
		Tenant: tnt,
	}

	t, err := c.Tenant().Create(ctx, connect.NewRequest(tcr))
	if err != nil {
		return err
	}
	log.Info("created tenant", "tenant", t)

	// try to create the same tenant with the returned id another time...
	tcr.Tenant.Meta = t.Msg.Tenant.Meta
	_, err = c.Tenant().Create(ctx, connect.NewRequest(tcr))

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
	tfrs, err := c.Tenant().Find(ctx, connect.NewRequest(tfrq))
	if err != nil {
		return err
	}
	for i := range tfrs.Msg.Tenants {
		log.Info("found tenant", "tenant", tfrs.Msg.Tenants[i])
	}

	log.Info("get tenant with id")
	tgr := &v1.TenantGetRequest{
		Id: t.Msg.Tenant.Meta.Id,
	}
	tgres, err := c.Tenant().Get(ctx, connect.NewRequest(tgr))
	if err != nil {
		return err
	}
	log.Info("got tenant", "id", tgres)

	v1t := mapper.ToV1Tenant(tgres.Msg.Tenant)
	mdm1t := mapper.ToMdmV1Tenant(v1t)

	_, err = c.Tenant().Update(ctx, connect.NewRequest(&v1.TenantUpdateRequest{Tenant: mdm1t}))
	if err != nil {
		return err

	}

	log.Info("get tenant with non-existent id")
	tgrNotFound := &v1.TenantGetRequest{
		Id: "1982739817298219873",
	}
	_, err = c.Tenant().Get(ctx, connect.NewRequest(tgrNotFound))
	if !v1.IsNotFound(err) {
		return err
	}

	// get tenant one more time to have some older version after update to provoke an optimistic lock error
	tgres2, _ := c.Tenant().Get(ctx, connect.NewRequest(tgr))

	tgres, err = c.Tenant().Get(ctx, connect.NewRequest(tgr))
	if err != nil {
		return err
	}
	tenant := tgres.Msg.Tenant
	tenant.Name = "some other name"

	tur := &v1.TenantUpdateRequest{
		Tenant: tenant,
	}
	tures, err := c.Tenant().Update(ctx, connect.NewRequest(tur))
	if err != nil {
		return err
	}
	log.Info("updated tenant", "tenant", tures)

	tenant2 := tgres2.Msg.Tenant
	tenant2.Name = "update older tenant"
	_, err = c.Tenant().Update(ctx, connect.NewRequest(tur))
	if !v1.IsOptimistickLockError(err) {
		return err
	}

	log.Info("find tenant with id")
	tfrqi := &v1.TenantFindRequest{
		Id: pointer.Pointer(t.Msg.Tenant.Meta.Id),
	}
	tfrsi, err := c.Tenant().Find(ctx, connect.NewRequest(tfrqi))
	if err != nil {
		return err
	}
	for i := range tfrsi.Msg.Tenants {
		log.Info("found tenant", "tenant", tfrsi.Msg.Tenants[i])
	}

	log.Info("delete tenant with id")
	tdr := &v1.TenantDeleteRequest{
		Id: t.Msg.Tenant.Meta.Id,
	}
	_, err = c.Tenant().Delete(ctx, connect.NewRequest(tdr))
	if err != nil {
		return err
	}

	log.Info("try to delete already deleted tenant")
	tdr2 := &v1.TenantDeleteRequest{
		Id: t.Msg.Tenant.Meta.Id,
	}
	_, err = c.Tenant().Delete(ctx, connect.NewRequest(tdr2))
	if !v1.IsNotFound(err) {
		log.Info("got expected grpc code, indicating not found")
	}

	pbHt := timestamppb.Now()
	thr, err := c.Tenant().GetHistory(ctx, connect.NewRequest(&v1.TenantGetHistoryRequest{
		Id: tdr.Id,
		At: pbHt,
	}))
	if err != nil {
		return err
	}
	if thr.Msg.Tenant.Name != "some other name" {
		return err
	}
	log.Info("found history tenant", "tenant", thr.Msg.Tenant)
	return nil
}
