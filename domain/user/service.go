package user

// UserService handles the business logic of the User module
type UserService struct {
	// Add any dependencies here
	repo *UserRepository
}

// NewUserService creates a new instance of TestService
func NewUserService(repo *UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// GetMessage returns a greeting message
func (s *UserService) GetMessage() UserModel {
	return s.repo.GetMessage()
}
