package services

// User defines FOXDEN user
type User struct {
	Name   string
	Groups []string
	Scopes []string
}

// UserAttributes defines generic interface of Foxden user attributes
type UserAttributes interface {
	// GetUsers should return all Foxden user names
	GetUsers() ([]string, error)
	// GetGroups should return all existing groups
	GetGroups() ([]string, error)
	// Get should return User
	Get(user string) (User, error)
}
