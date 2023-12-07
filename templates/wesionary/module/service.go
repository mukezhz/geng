package {{.PackageName}}

// {{.ModuleName}}Service handles the business logic of the {{.ModuleName}} module
type {{.ModuleName}}Service struct {
	// Add any dependencies here
	repo *{{.ModuleName}}Repository
}

// New{{.ModuleName}}Service creates a new instance of TestService
func New{{.ModuleName}}Service(repo *{{.ModuleName}}Repository) *{{.ModuleName}}Service {
	return &{{.ModuleName}}Service{
		repo: repo,
	}
}

// GetMessage returns a greeting message
func (s *{{.ModuleName}}Service) GetMessage() {{.ModuleName}}Model {
	return s.repo.GetMessage()
}
