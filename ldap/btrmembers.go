package ldap

import (
	"fmt"
	"log"
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

// BtrMembersUids provides list of CLASSE users which associated with BTR
func BtrMembersUids(login, password, btr string) ([]string, error) {
	var records []string
	var err error

	// extract LDAP configuration parameters
	ldapURL := srvConfig.Config.LDAP.URL
	baseDN := srvConfig.Config.LDAP.BaseDN
	attributes := []string{"member"}
	results, err := SearchBy(ldapURL, login, password, baseDN, btr, "cn", attributes)
	if err != nil {
		return records, fmt.Errorf("[golib.ldap.BtrMembersUids] Search error: %w", err)
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
	// now search for uids for each member
	var out []string
	for _, m := range records {
		attributes := []string{"uid"}
		results, err := SearchBy(ldapURL, login, password, baseDN, m, "cn", attributes)
		if err != nil {
			log.Printf("[golib.ldap.BtrMembersUids] Search error: %v", err)
			continue
		}
		for _, entry := range results.Entries {
			uid := entry.GetAttributeValue("uid")
			rec := fmt.Sprintf("%s uid %s", m, uid)
			out = append(out, rec)
		}
	}
	if len(out) == len(records) {
		return out, nil
	}
	return records, nil
}
