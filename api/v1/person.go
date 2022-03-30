package v1

//go:generate go run ../../pkg/gen/genscanvaluer.go -package v1 -type Person

func (m *Person) NewPersonResponse() *PersonResponse {
	return &PersonResponse{
		Person: m,
	}
}

func (m *PersonDeleteRequest) NewPerson() *Person {
	return &Person{
		Meta: &Meta{Id: m.Id},
	}
}
