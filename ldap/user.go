package ldap

import (
	"fmt"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
)

// UserInfo represents the structure returned by the user service.
type UserInfo struct {
	DN        string    `json:"DN"`
	Name      string    `json:"Name"`
	Email     string    `json:"Email"`
	Uid       string    `json:"Uid"`
	UidNumber int       `json:"UidNumber"`
	GidNumber int       `json:"GidNumber"`
	Groups    []string  `json:"Groups"`
	Btrs      []string  `json:"Btrs"`
	Beamlines []string  `json:"Beamlines"`
	Expire    time.Time `json:"Expire"`
	Foxdens   []string  `json:"Foxdens"`
}

// GetEmail lookup user's email for provided user name
func GetEmail(login, password, name string) (string, error) {
	var email string

	// extract LDAP configuration parameters
	ldapURL := srvConfig.Config.LDAP.URL
	baseDN := srvConfig.Config.LDAP.BaseDN
	attributes := []string{"mail"}
	results, err := SearchBy(
		ldapURL, login, password, baseDN, name, "cn", attributes)
	if err != nil {
		return email, fmt.Errorf("[golib.ldap.GetEmail] Search error: %w", err)
	}
	for _, entry := range results.Entries {
		email = entry.GetAttributeValue("mail")
		return email, nil
	}
	return email, nil
}
