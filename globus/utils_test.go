package globus

import (
	"fmt"
	"net/url"
	"testing"
)

// TestGlobusLink
func TestGlobusLink(t *testing.T) {
	cid := "collection-id"
	path := "/nfs/chess/aux/cycles/2024-3/id1a3/ko-3538-d/raw_data/align-0925/25"
	gurl, err := GlobusLink(cid, path)
	if err != nil {
		t.Error(err)
	}
	spath := url.QueryEscape("/2024-3/id1a3/ko-3538-d/align-0925")
	epath := fmt.Sprintf("https://app.globus.org/file-manager?origin_id=%s&origin_path=%s", cid, spath)
	if gurl != epath {
		msg := fmt.Sprintf("ERROR: got\n%s\nexpect:\n%s", gurl, epath)
		t.Error(msg)
	}

	// test /nfs/chess/raw/2025-1/id3b/thompson-4135-a raw path
	path = "/nfs/chess/raw/2025-1/id3b/thompson-4135-a"
	gurl, err = GlobusLink(cid, path)
	if err != nil {
		t.Error(err)
	}
	spath = url.QueryEscape("/2025-1/id3b/thompson-4135-a")
	epath = fmt.Sprintf("https://app.globus.org/file-manager?origin_id=%s&origin_path=%s", cid, spath)
	if gurl != epath {
		msg := fmt.Sprintf("ERROR: got\n%s\nexpect:\n%s", gurl, epath)
		t.Error(msg)
	}

	// test for empty path
	_, err = GlobusLink(cid, "")
	if err == nil {
		t.Error("do not get error for empty path")
	}

	// test non chess path
	_, err = GlobusLink(cid, "/bla")
	if err == nil {
		t.Error("do not get error for non chess path test")
	}
}
