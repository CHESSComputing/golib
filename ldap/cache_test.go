package ldap

import (
	"testing"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/stretchr/testify/assert"
)

// Mock configuration initialization
func mockConfig(expire int) {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	srvConfig.Config.LDAP.Expire = expire
	//	srvConfig.Config = &srvConfig.Configuration{
	//	    LDAP: srvConfig.LDAPConfig{
	//	        Expire: expire, // Expiration time in seconds
	//	    },
	//	}
}

// Test for Cache Entry Expiration
func TestCacheEntryExpiration(t *testing.T) {
	mockConfig(2) // Set cache expiration to 2 seconds
	cache := Cache{Map: make(map[string]Entry)}

	entry := Entry{
		DN:     "cn=test,ou=users,dc=example,dc=com",
		Groups: []string{"CN=Users"},
		Expire: time.Now().Add(1 * time.Second),
	}

	// Store entry in cache
	cache.Map["testuser"] = entry

	// Initially, the entry should be valid
	retrievedEntry, err := cache.Search("testuser", "testpassword", "testuser")
	assert.Nil(t, err, "Expected no error when fetching valid cache entry")
	assert.Equal(t, entry.DN, retrievedEntry.DN, "Expected same DN from cache")

	// Wait for expiration
	time.Sleep(2 * time.Second)

	// Now, the entry should be expired
	_, err = cache.Search("testuser", "testpassword", "testuser")
	assert.NotNil(t, err, "Expected error when fetching expired cache entry")
}

// Mock LDAP Search function
/*
func mockSearch(ldapURL, login, password, baseDN, user string, attributes []string) ([]Entry, error) {
	entry := Entry{
		DN:     "cn=test,ou=users,dc=example,dc=com",
		Groups: []string{"CN=Users"},
		Expire: time.Now().Add(5 * time.Second),
	}
	return []Entry{entry}, nil
}
*/

// Test LDAP Search Functionality
func TestLDAPSearch(t *testing.T) {
	mockConfig(5) // Set expiration time to 5 seconds
	cache := Cache{Map: make(map[string]Entry)}

	// Simulate a fresh LDAP query
	entry := Entry{
		DN:     "cn=test,ou=users,dc=example,dc=com",
		Groups: []string{"CN=Users"},
		Expire: time.Now().Add(5 * time.Second),
	}
	cache.Map["testuser"] = entry

	retrievedEntry, err := cache.Search("testuser", "testpassword", "testuser")
	assert.Nil(t, err, "Expected successful cache retrieval")
	assert.Equal(t, "cn=test,ou=users,dc=example,dc=com", retrievedEntry.DN, "DN should match the expected value")
}

// Test Cache Expiration and Refresh
func TestCacheExpirationAndRefresh(t *testing.T) {
	mockConfig(2) // Set cache expiration to 2 seconds
	cache := Cache{Map: make(map[string]Entry)}

	// Insert an entry that expires in 2 seconds
	entry := Entry{
		DN:     "cn=expired,ou=users,dc=example,dc=com",
		Groups: []string{"CN=Admins"},
		Expire: time.Now().Add(1 * time.Second),
	}
	cache.Map["expireduser"] = entry

	// Verify it's retrievable
	retrievedEntry, err := cache.Search("testuser", "testpassword", "expireduser")
	assert.Nil(t, err)
	assert.Equal(t, "cn=expired,ou=users,dc=example,dc=com", retrievedEntry.DN)

	// Wait for expiration
	time.Sleep(2 * time.Second)

	// It should now return an error since entry is expired
	_, err = cache.Search("testuser", "testpassword", "expireduser")
	assert.NotNil(t, err, "Expected error when retrieving expired cache entry")
}
