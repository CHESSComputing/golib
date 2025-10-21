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

// helper function to handle /nfs/chess/aux/cycles paths
func auxCyclesPath(path string) (string, error) {
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
	return epath, nil
}

// helper function to handle /nfs/chess/raw paths
func rawPath(path string) (string, error) {
	epath := strings.Replace(path, "/nfs/chess/raw", "", -1)
	epath = url.QueryEscape(epath)
	return epath, nil
}

// generate Globus collection URL link
func GlobusLink(cid, path string) (string, error) {
	if path == "" {
		return "", errors.New("empty path")
	}
	var epath string
	var err error
	if strings.HasPrefix(path, "/nfs/chess/aux/cycles") {
		epath, err = auxCyclesPath(path)
	} else if strings.HasPrefix(path, "/nfs/chess/raw/") {
		epath, err = rawPath(path)
	} else {
		return "", errors.New(fmt.Sprintf("Non chess path %s", path))
	}
	if err != nil {
		return "", err
	}
	gurl := fmt.Sprintf("https://app.globus.org/file-manager?origin_id=%s&origin_path=%s", cid, epath)
	return gurl, nil
}

// ChessGlobusLink provides globus link to given globus collection name and path
func ChessGlobusLink(collection, path string) (string, error) {
	if globusCache == nil {
		globusCache = make(map[string]string)
	}

	// if FOXDEN configuration provides Globus OriginID we will use it as collection id
	if srvConfig.Config.Globus.OriginID != "" {
		return GlobusLink(srvConfig.Config.Globus.OriginID, path)
	}

	var cid string
	// check if globusCache has our collection
	mu := sync.Mutex{}
	mu.Lock()
	cid, ok := globusCache[collection]
	mu.Unlock()
	if ok {
		return GlobusLink(cid, path)
	}

	// obtain collection id from globus
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
