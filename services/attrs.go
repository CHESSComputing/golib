package services

// User defines FOXDEN user
type User struct {
	Name         string
	Groups       []string
	Scopes       []string
	Btrs         []string
	FoxdenGroups []string
}

// UserAttributes defines generic interface of Foxden user attributes
type UserAttributes interface {
	// Init function initialize foxden user
	Init()
	// GetUsers should return all Foxden user names
	GetUsers() ([]string, error)
	// GetGroups should return all existing groups
	GetGroups() ([]string, error)
	// Get should return User
	Get(user string) (User, error)
	// GetGroup should return group associated with given did
	GetGroup(did string) string
	// GetEmail should return email associated with this user
	GetEmail(user string) (string, error)
	// GetMembers should return group members associated with this user
	GetMembers(user string) ([]string, error)
}
