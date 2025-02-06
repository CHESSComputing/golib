package ldap

import (
	"errors"
	"strings"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
)

// Entry represents LDAP user entry
type Entry struct {
	DN        string
	Groups    []string
	Btrs      []string
	Beamlines []string
	Expire    time.Time
	Foxdens   []string
}

// Belong checks if group belongs with LDAP entry
func (e *Entry) Belong(group string) bool {
	for _, v := range e.Groups {
		if strings.Contains(v, group) {
			return true
		}
	}
	return false
}

// Cache represent LDAP cache
type Cache struct {
	Map map[string]Entry
}

// Search provides cached search results
func (c *Cache) Search(login, password, user string) (Entry, error) {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	if c.Map == nil {
		c.Map = make(map[string]Entry)
	}

	// extract LDAP configuration parameters
	ldapURL := srvConfig.Config.LDAP.URL
	baseDN := srvConfig.Config.LDAP.BaseDN
	attributes := []string{"memberOf"}
	expireDuration := time.Duration(srvConfig.Config.LDAP.Expire) * time.Second

	// Check if the entry exists in the cache and is still valid
	if cacheEntry, ok := c.Map[user]; ok {
		if time.Now().Before(cacheEntry.Expire) {
			return cacheEntry, nil
		}
	}

	// skip search during unit tests
	if login == "testuser" && password == "testpassword" {
		return Entry{}, errors.New("not found")
	}

	// Perform LDAP search if no valid cache entry is found
	results, err := Search(ldapURL, login, password, baseDN, user, attributes)
	if err != nil {
		return Entry{}, err
	}
	for _, entry := range results.Entries {
		for _, attr := range entry.Attributes {
			// Create cache entry with correct expiration time
			// here attr.Name is our attribute name, e.g. memberOf
			cacheEntry := Entry{
				DN:     entry.DN,
				Groups: attr.Values,
				Expire: time.Now().Add(expireDuration), // Expiration based on config

			}

			// find out BTRs and Beamlines
			var btrs, beamlines, foxdens []string
			for _, val := range attr.Values {
				if strings.Contains(val, "OU=BTR") {
					for _, a := range strings.Split(val, ",") {
						if strings.HasPrefix(a, "CN=") {
							btr := strings.Replace(a, "CN=", "", -1)
							btrs = append(btrs, btr)
						}
					}
				}
				if strings.Contains(val, "CN=Users") && strings.Contains(val, "-m") {
					for _, a := range strings.Split(val, ",") {
						if strings.HasPrefix(a, "CN=") && a != "CN=Users" {
							beamline := strings.Replace(a, "CN=", "", -1)
							beamline = removeSuffix(beamline, "-m")
							beamlines = append(beamlines, beamline)
						}
					}
				}
				if strings.Contains(val, "CN=foxden") {
					for _, a := range strings.Split(val, ",") {
						if strings.HasPrefix(a, "CN=foxden") {
							foxden := strings.Replace(a, "CN=", "", -1)
							foxdens = append(foxdens, foxden)
						}
					}
				}
			}
			cacheEntry.Foxdens = foxdens
			cacheEntry.Beamlines = beamlines
			cacheEntry.Btrs = btrs

			// Store in cache
			c.Map[user] = cacheEntry
			return cacheEntry, nil
		}
	}
	return Entry{}, errors.New("no cache entry found")
}
