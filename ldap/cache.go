package ldap

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/CHESSComputing/golib/utils"
)

// Entry represents LDAP user entry
type Entry struct {
	DN        string
	Name      string
	Email     string
	Uid       string
	UidNumber int
	GidNumber int
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
	mutex   sync.RWMutex
	Map     map[string]Entry
	Verbose int
}

// Search provides cached search results
func (c *Cache) Search(login, password, user string) (Entry, error) {
	// by default we search by uid
	return c.SearchBy(login, password, user, "uid")
}

// SearchBy provides cached search results
func (c *Cache) SearchBy(login, password, user, method string) (Entry, error) {
	if c.Map == nil {
		c.Map = make(map[string]Entry)
	}

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
	records, err := Records(
		login, password, user, method, c.Verbose)
	if err != nil {
		return Entry{}, fmt.Errorf("[golib.ldap.Cache.Search] Search error: %w", err)
	}
	if len(records) > 0 && err == nil {
		// Store in cache
		c.mutex.Lock()
		c.Map[user] = cacheEntry
		c.mutex.Unlock()
		return cacheEntry, nil
	}
	return Entry{}, errors.New("no cache entry found")
}

// MultiSearch provides cached search results
func Records(login, password, user, method string, verbose int) ([]Entry, error) {
	var records []Entry
	var err error

	// extract LDAP configuration parameters
	ldapURL := srvConfig.Config.LDAP.URL
	baseDN := srvConfig.Config.LDAP.BaseDN
	attributes := []string{
		"memberOf", "uidNumber", "gidNumber", "uid", "mail", "displayName"}
	expireDuration := time.Duration(srvConfig.Config.LDAP.Expire) * time.Second

	// Perform LDAP search if no valid cache entry is found
	results, err := SearchBy(
		ldapURL, login, password, baseDN, user, method, attributes)
	if err != nil {
		return records, fmt.Errorf("[golib.ldap.Cache.Search] Search error: %w", err)
	}
	for _, entry := range results.Entries {
		// Create cache entry with correct expiration time
		cacheEntry := Entry{
			DN:     entry.DN,
			Expire: time.Now().Add(expireDuration), // Expiration based on config

		}
		cacheEntry.Name = entry.GetAttributeValue("displayName")
		cacheEntry.Uid = entry.GetAttributeValue("uid")
		cacheEntry.Email = entry.GetAttributeValue("mail")
		cacheEntry.Name = entry.GetAttributeValue("displayName")
		if uid, err := strconv.ParseInt(entry.GetAttributeValue("uidNumber"), 10, 32); err == nil {
			cacheEntry.UidNumber = int(uid)
		}
		if gid, err := strconv.ParseInt(entry.GetAttributeValue("gidNumber"), 10, 32); err == nil {
			cacheEntry.GidNumber = int(gid)
		}
		// work with memberOf attributes
		attrValues := entry.GetAttributeValues("memberOf")
		cacheEntry.Groups = attrValues

		// find out BTRs and Beamlines
		var btrs, beamlines, foxdens []string
		for _, val := range attrValues {
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
						if verbose > 0 {
							log.Printf("### recursive search groupCN=%s, users=%+v, error=%v", groupCN, users, err)
						}
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

		records = append(records, cacheEntry)
	}
	return records, err
}
