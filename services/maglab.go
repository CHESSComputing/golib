package services

import (
	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/CHESSComputing/golib/ldap"
)

// MaglabUser represents Maglab user with UserAttributes interface APIs
type MaglabUser struct {
	URL       string // ldap URL
	Login     string // ldap user name
	Password  string // ldap password
	BaseDN    string // ldap baseDN
	ldapCache *ldap.Cache
}

// Init function initialize MaglabUser ldap cache
func (c *MaglabUser) Init() {
	if c.URL == "" {
		c.URL = srvConfig.Config.LDAP.URL

	}
	if c.Login == "" {
		c.Login = srvConfig.Config.LDAP.Login

	}
	if c.Password == "" {
		c.Password = srvConfig.Config.LDAP.Password

	}
	if c.BaseDN == "" {
		c.BaseDN = srvConfig.Config.LDAP.BaseDN

	}
	c.ldapCache = &ldap.Cache{Map: make(map[string]ldap.Entry)}
}

// GetUsers implement UserAttributes GetUsers API
func (c *MaglabUser) GetUsers() ([]string, error) {
	return ldap.GetUsers(c.URL, c.Login, c.Password, c.BaseDN)
}

// GetGroups implements UserAttributes Get API
func (c *MaglabUser) GetGroups() ([]string, error) {
	return ldap.GetGroups(c.URL, c.Login, c.Password, c.BaseDN)
}

// Get implements UserAttributes Get API
func (c *MaglabUser) Get(name string) (User, error) {
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
