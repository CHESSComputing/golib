package ldap

import (
	"crypto/tls"
	"fmt"
	"log"
	"strings"
	"time"

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

func isBTRUserDN(dn string) bool {
	return strings.Contains(strings.ToUpper(dn), "OU=BTR")
}

func getGroupMembersRecursive(
	l *ldap.Conn,
	groupDN string,
	baseDN string,
	visited map[string]bool,
	results map[string]struct{},
	recursionLevel, verbose int,
) error {

	// prevent infinite loops
	if visited[groupDN] {
		return nil
	}
	visited[groupDN] = true

	req := ldap.NewSearchRequest(
		groupDN,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=*)",
		[]string{"memberOf"},
		nil,
	)

	res, err := l.Search(req)
	if err != nil {
		return err
	}

	if len(res.Entries) == 0 {
		return nil
	}

	for _, memberDN := range res.Entries[0].GetAttributeValues("memberOf") {

		// Case 1: member is a user
		if isBTRUserDN(memberDN) {
			results[memberDN] = struct{}{}
			continue
		}

		// Case 2: member might be another group
		memberReq := ldap.NewSearchRequest(
			memberDN,
			ldap.ScopeBaseObject,
			ldap.NeverDerefAliases,
			0, 0, false,
			"(objectClass=posixGroup)",
			[]string{"cn"},
			nil,
		)

		memberRes, err := l.Search(memberReq)
		if err == nil && len(memberRes.Entries) > 0 {
			// recurse into nested group
			if recursionLevel <= 1 {
				// reached recursion limit, skip further recursion
				return nil
			}
			_ = getGroupMembersRecursive(l, memberDN, baseDN, visited, results, recursionLevel-1, verbose)
		}
	}

	if verbose > 1 {
		log.Printf("INFO: recursive BTR lookup stoped at recursion level=%d", recursionLevel)
	}

	return nil
}

func GetBTRUsersFromGroup(
	ldapURL, login, password, baseDN, groupCN string, recursionLevel, verbose int,
) ([]string, error) {

	time0 := time.Now()
	defer func() {
		log.Printf("INFO: GetBTRUsersFromGroup with recursion level=%d elapsed time: %s", recursionLevel, time.Since(time0))
	}()
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	l, err := ldap.DialURL(ldapURL, ldap.DialWithTLSConfig(tlsConfig))
	if err != nil {
		return nil, err
	}
	defer l.Close()

	if err := l.Bind(login, password); err != nil {
		return nil, err
	}

	groupDN := fmt.Sprintf("CN=%s,CN=Users,%s", groupCN, baseDN)

	visited := make(map[string]bool)
	results := make(map[string]struct{})

	if err := getGroupMembersRecursive(l, groupDN, baseDN, visited, results, recursionLevel, verbose); err != nil {
		return nil, err
	}

	var users []string
	for dn := range results {
		users = append(users, dn)
	}

	return users, nil
}

// GetBTR extracts BTR value from given user DN value
func GetBTR(val string) string {
	if strings.Contains(val, "OU=BTR") {
		for _, a := range strings.Split(val, ",") {
			if strings.HasPrefix(a, "CN=") {
				btr := strings.Replace(a, "CN=", "", -1)
				return btr
			}
		}
	}
	return ""
}
