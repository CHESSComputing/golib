package ldap

import (
	"fmt"
	"strings"

	srvConfig "github.com/CHESSComputing/golib/config"
)

// BtrMembers provides list of CLASSE users which associated with BTR
func BtrMembers(login, password, btr string) ([]string, error) {
	var records []string
	var err error

	// extract LDAP configuration parameters
	ldapURL := srvConfig.Config.LDAP.URL
	baseDN := srvConfig.Config.LDAP.BaseDN
	attributes := []string{"member"}
	results, err := SearchBy(
		ldapURL, login, password, baseDN, btr, "cn", attributes)
	if err != nil {
		return records, fmt.Errorf("[golib.ldap.BtrMembers] Search error: %w", err)
	}
	for _, entry := range results.Entries {
		for _, cn := range entry.GetAttributeValues("member") {
			if cn != "" {
				arr := strings.Split(cn, ",")
				name := strings.Replace(arr[0], "CN=", "", -1)
				records = append(records, name)
			}
		}
	}
	return records, nil
}
