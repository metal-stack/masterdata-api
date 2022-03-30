package v1

//go:generate go run ../../pkg/gen/genscanvaluer.go -package v1 -type User

func (m *User) NewUserResponse() *UserResponse {
	return &UserResponse{
		User: m,
	}
}

func (m *UserDeleteRequest) NewUser() *User {
	return &User{
		Meta: &Meta{Id: m.Id},
	}
}
