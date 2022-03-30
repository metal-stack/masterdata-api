package service

import (
	"context"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/datastore"
	"go.uber.org/zap"
)

type UserService struct {
	Storage datastore.Storage
	log     *zap.Logger
}

func NewUserService(s datastore.Storage, l *zap.Logger) *UserService {
	return &UserService{
		Storage: NewStorageStatusWrapper(s),
		log:     l,
	}
}

func (s *UserService) Create(ctx context.Context, req *v1.UserCreateRequest) (*v1.UserResponse, error) {
	user := req.User

	// allow create without sending Meta
	if user.Meta == nil {
		user.Meta = &v1.Meta{}
	}
	err := s.Storage.Create(ctx, user)
	return user.NewUserResponse(), err
}
func (s *UserService) Update(ctx context.Context, req *v1.UserUpdateRequest) (*v1.UserResponse, error) {
	user := req.User
	err := s.Storage.Update(ctx, user)
	return user.NewUserResponse(), err
}
func (s *UserService) Delete(ctx context.Context, req *v1.UserDeleteRequest) (*v1.UserResponse, error) {
	user := req.NewUser()
	err := s.Storage.Delete(ctx, user)
	return user.NewUserResponse(), err
}
func (s *UserService) Get(ctx context.Context, req *v1.UserGetRequest) (*v1.UserResponse, error) {
	user := &v1.User{}
	err := s.Storage.Get(ctx, req.Id, user)
	if err != nil {
		return nil, err
	}
	return user.NewUserResponse(), nil
}
