package ldap

import (
	"crypto/tls"
	"fmt"
	"log"

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
