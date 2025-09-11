package services

import "github.com/CHESSComputing/golib/ldap"

// CHESSUser represents chess user with UserAttributes interface APIs
type CHESSUser struct {
	URL       string // ldap URL
	Login     string // ldap user name
	Password  string // ldap password
	BaseDN    string // ldap baseDN
	ldapCache *ldap.Cache
}

// Init function initialize CHESSUser ldap cache
func (c *CHESSUser) Init() {
	c.ldapCache = &ldap.Cache{Map: make(map[string]ldap.Entry)}
}

// GetUsers implement UserAttributes GetUsers API
func (c *CHESSUser) GetUsers() ([]string, error) {
	return ldap.GetUsers(c.URL, c.Login, c.Password, c.BaseDN)
}

// GetGroups implements UserAttributes Get API
func (c *CHESSUser) GetGroups() ([]string, error) {
	return ldap.GetUsers(c.URL, c.Login, c.Password, c.BaseDN)
}

// Get implements UserAttributes Get API
func (c *CHESSUser) Get(name string) (User, error) {
	user := User{
		Name: name,
	}
	entry, err := c.ldapCache.Search(c.Login, c.Password, name)
	if err != nil {
		return user, err
	}
	user.Groups = entry.Btrs
	return user, nil
}
