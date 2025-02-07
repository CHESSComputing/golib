package globus

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	srvConfig "github.com/CHESSComputing/golib/config"
)

// keep map of globus collection vs ids for caching purposes
var globusCache map[string]string

// generate Globus collection URL link
func GlobusLink(cid, path string) (string, error) {
	if path == "" {
		return "", errors.New("empty path")
	}
	if !strings.HasPrefix(path, "/nfs/chess/aux/cycles") {
		return "", errors.New(fmt.Sprintf("Non chess path %s", path))
	}
	// parse the NFS path and extract relevant parts
	path = strings.Replace(path, "/nfs/chess/aux/cycles", "", -1)
	arr := strings.Split(path, "/raw_data/")
	if len(arr) != 2 {
		msg := fmt.Sprintf("Unable to parse %s", path)
		return "", errors.New(msg)
	}
	idir := arr[0]                        // leading part
	edir := strings.Split(arr[1], "/")[0] // ending part (first directory)
	epath := url.QueryEscape(fmt.Sprintf("%s/%s", idir, edir))
	gurl := fmt.Sprintf("https://app.globus.org/file-manager?origin_id=%s&&origin_path=%s", cid, epath)
	return gurl, nil
}

// ChessGlobusLink provides globus link to given globus collection name and path
func ChessGlobusLink(collection, path string) (string, error) {
	if globusCache == nil {
		globusCache = make(map[string]string)
	}

	// if FOXDEN configuration provides Globus OriginID we will use it as collection id
	if srvConfig.Config.Globus.OriginID != "" {
		return GlobusLink(srvConfig.Config.Globus.ClientID, path)
	}

	// check if globusCache has our collection
	if cid, ok := globusCache[collection]; ok {
		return GlobusLink(cid, path)
	}

	// obtain collection id from globus
	var cid string
	scopes := []string{"urn:globus:auth:scope:transfer.api.globus.org:all"}
	token, err := Token(scopes)
	if err != nil {
		return "", err
	}
	// find collection id
	mapMutex := sync.RWMutex{}
	records := Search(token, collection)
	for _, r := range records {
		if r.Name == collection {
			cid = r.Id
			mapMutex.Lock()
			globusCache[collection] = cid
			mapMutex.Unlock()
			break
		}
	}
	if cid == "" {
		return "", errors.New(fmt.Sprintf("No Globus collection found for collection %s", collection))
	}
	return GlobusLink(cid, path)
}
