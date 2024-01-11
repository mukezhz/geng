package user

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (s *UserRepository) GetMessage() UserModel {
	return UserModel{Message: "Hello World"}
}
