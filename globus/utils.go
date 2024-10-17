package globus

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

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
	var cid string
	scopes := []string{"urn:globus:auth:scope:transfer.api.globus.org:all"}
	token, err := Token(scopes)
	if err != nil {
		return "", err
	}
	// find collection id
	records := Search(token, collection)
	for _, r := range records {
		if r.Name == collection {
			cid = r.Id
			break
		}
	}
	if cid == "" {
		return "", errors.New(fmt.Sprintf("No Globus collection found for collection %s", collection))
	}
	return GlobusLink(cid, path)
}
