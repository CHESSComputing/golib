package ldap

import (
	"crypto/tls"
	"fmt"
	"log"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// Search provides LDAP search for our FOXDEN LDAP service
func Search(ldapURL, login, password, baseDN, user string, attributes []string) (*ldap.SearchResult, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	l, err := ldap.DialURL(ldapURL, ldap.DialWithTLSConfig(tlsConfig))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// connect to LDAP server
	err = l.Bind(login, password)
	if err != nil {
		log.Fatal(err)
	}

	// create LDAP filter, must start and finish with ()!
	//     filter := fmt.Sprintf("(CN=%s)", ldap.EscapeFilter(user))
	filter := fmt.Sprintf("(uid=%s)", ldap.EscapeFilter(user))

	// query LDAP server
	sizeLimit := 0
	timeLimit := 0
	typesOnly := false
	controls := []ldap.Control{}
	searchReq := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		sizeLimit,
		timeLimit,
		typesOnly,
		filter,
		attributes,
		controls)
	result, err := l.Search(searchReq)
	return result, err
}

// helper function to remove specific suffix from the end of the string
func removeSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

// GetUsers queries LDAP and returns all user CNs
func GetUsers(ldapURL, login, password, baseDN string) ([]string, error) {
	var users []string
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	l, err := ldap.DialURL(ldapURL, ldap.DialWithTLSConfig(tlsConfig))
	if err != nil {
		return users, err
	}
	defer l.Close()

	// connect to LDAP server
	err = l.Bind(login, password)
	if err != nil {
		return users, err
	}

	// Adjust filter depending on your LDAP schema
	// For Active Directory: "(objectClass=user)"
	// For OpenLDAP: "(objectClass=inetOrgPerson)"
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=inetOrgPerson)",
		[]string{"cn"},
		nil,
	)

	result, err := l.Search(searchRequest)
	if err != nil {
		return users, err
	}

	for _, entry := range result.Entries {
		users = append(users, entry.GetAttributeValue("cn"))
	}
	return users, nil
}

// GetGroups queries LDAP and returns all group CNs
func GetGroups(ldapURL, login, password, baseDN string) ([]string, error) {
	var groups []string
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	l, err := ldap.DialURL(ldapURL, ldap.DialWithTLSConfig(tlsConfig))
	if err != nil {
		return groups, err
	}
	defer l.Close()

	// connect to LDAP server
	err = l.Bind(login, password)
	if err != nil {
		log.Fatal(err)
	}

	// For Active Directory: "(objectClass=group)"
	// For OpenLDAP: "(objectClass=posixGroup)"
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=posixGroup)",
		[]string{"cn"},
		nil,
	)

	result, err := l.Search(searchRequest)
	if err != nil {
		return groups, err
	}

	for _, entry := range result.Entries {
		groups = append(groups, entry.GetAttributeValue("cn"))
	}
	return groups, nil
}
