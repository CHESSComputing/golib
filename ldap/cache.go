package ldap

import (
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/CHESSComputing/golib/utils"
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
	mutex sync.RWMutex
	Map   map[string]Entry
}

// Search provides cached search results
func (c *Cache) Search(login, password, user string) (Entry, error) {
	if c.Map == nil {
		c.Map = make(map[string]Entry)
	}

	// extract LDAP configuration parameters
	ldapURL := srvConfig.Config.LDAP.URL
	baseDN := srvConfig.Config.LDAP.BaseDN
	attributes := []string{"memberOf"}
	expireDuration := time.Duration(srvConfig.Config.LDAP.Expire) * time.Second

	// Check if the entry exists in the cache and is still valid
	c.mutex.RLock()
	cacheEntry, ok := c.Map[user]
	c.mutex.RUnlock()
	if ok {
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
				btr := GetBTR(val)
				if btr != "" {
					btrs = append(btrs, btr)
				}
				// perform recursive search of users groups to extract additional BTRs
				if strings.Contains(val, "CN=Users") {
					arr := strings.Split(val, ",")
					if len(arr) > 0 {
						a := arr[0] // first CN attribute represents group
						groupCN := strings.Replace(a, "CN=", "", -1)
						verbose := srvConfig.Config.Authz.Verbose
						recursionLevel := srvConfig.Config.LDAP.RecursionLevel
						if recursionLevel == 0 {
							recursionLevel = 5 // default value based on CLASSE IT suggestion
						}
						users, err := GetBTRUsersFromGroup(ldapURL, login, password, baseDN, groupCN, recursionLevel, verbose)
						if err == nil {
							for _, user := range users {
								btr := GetBTR(user)
								if btr != "" {
									btrs = append(btrs, btr)
								}
							}
						} else {
							log.Printf("### recursive search groupCN=%s, users=%+v, error=%v", groupCN, users, err)
						}
					}
				}
				btrs = utils.List2Set(btrs)
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
			c.mutex.Lock()
			c.Map[user] = cacheEntry
			c.mutex.Unlock()
			return cacheEntry, nil
		}
	}
	return Entry{}, errors.New("no cache entry found")
}
