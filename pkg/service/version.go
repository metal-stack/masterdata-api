package service

import (
	"context"

	"connectrpc.com/connect"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/v"
)

type versionService struct {
}

func NewVersionService() *versionService {
	return &versionService{}
}
func (vs *versionService) Get(context.Context, *connect.Request[v1.GetVersionRequest]) (*connect.Response[v1.GetVersionResponse], error) {
	res := v1.GetVersionResponse{Version: v.Version, Revision: v.Revision, BuildDate: v.BuildDate, GitSha1: v.GitSHA1}
	return connect.NewResponse(&res), nil
}
