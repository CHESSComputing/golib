package services

import (
	"strings"

	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/CHESSComputing/golib/ldap"
	"github.com/CHESSComputing/golib/utils"
)

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
func (c *CHESSUser) GetUsers() ([]string, error) {
	return ldap.GetUsers(c.URL, c.Login, c.Password, c.BaseDN)
}

// GetGroups implements UserAttributes Get API
func (c *CHESSUser) GetGroups() ([]string, error) {
	return ldap.GetGroups(c.URL, c.Login, c.Password, c.BaseDN)
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
	var groups []string
	for _, rec := range entry.Groups {
		if strings.Contains(rec, "BTR") {
			// this is BTR entry and not user's group
			continue
		}
		for _, a := range strings.Split(rec, ",") {
			if strings.HasPrefix(a, "CN=") {
				grp := strings.Replace(a, "CN=", "", -1)
				groups = append(groups, grp)
			}
		}
	}
	user.Groups = utils.List2Set(groups)
	user.Btrs = entry.Btrs
	user.FoxdenGroups = entry.Foxdens
	// add default scope
	user.Scopes = append(user.Scopes, "read")
	// add more scopes based on user's groups
	for _, grp := range user.Groups {
		if grp == "foxdenadmin" {
			user.Scopes = append(user.Scopes, "delete")
		}
		if grp == "foxdenrw" {
			user.Scopes = append(user.Scopes, "write")
		}
	}
	return user, nil
}
