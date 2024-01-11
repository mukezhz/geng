package hello

// HelloService handles the business logic of the Hello module
type HelloService struct {
	// Add any dependencies here
	repo *HelloRepository
}

// NewHelloService creates a new instance of TestService
func NewHelloService(repo *HelloRepository) *HelloService {
	return &HelloService{
		repo: repo,
	}
}

// GetMessage returns a greeting message
func (s *HelloService) GetMessage() HelloModel {
	return s.repo.GetMessage()
}
