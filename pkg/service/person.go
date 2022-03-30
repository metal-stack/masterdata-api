package service

import (
	"context"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/datastore"
	"go.uber.org/zap"
)

type PersonService struct {
	Storage datastore.Storage
	log     *zap.Logger
}

func NewPersonService(s datastore.Storage, l *zap.Logger) *PersonService {
	return &PersonService{
		Storage: NewStorageStatusWrapper(s),
		log:     l,
	}
}

func (s *PersonService) Create(ctx context.Context, req *v1.PersonCreateRequest) (*v1.PersonResponse, error) {
	person := req.Person

	// allow create without sending Meta
	if person.Meta == nil {
		person.Meta = &v1.Meta{}
	}
	err := s.Storage.Create(ctx, person)
	return person.NewPersonResponse(), err
}
func (s *PersonService) Update(ctx context.Context, req *v1.PersonUpdateRequest) (*v1.PersonResponse, error) {
	person := req.Person
	err := s.Storage.Update(ctx, person)
	return person.NewPersonResponse(), err
}
func (s *PersonService) Delete(ctx context.Context, req *v1.PersonDeleteRequest) (*v1.PersonResponse, error) {
	person := req.NewPerson()
	err := s.Storage.Delete(ctx, person)
	return person.NewPersonResponse(), err
}
func (s *PersonService) Get(ctx context.Context, req *v1.PersonGetRequest) (*v1.PersonResponse, error) {
	person := &v1.Person{}
	err := s.Storage.Get(ctx, req.Id, person)
	if err != nil {
		return nil, err
	}
	return person.NewPersonResponse(), nil
}
