package services

import (
	"encoding/json"
	"fmt"
	"io"

	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/CHESSComputing/golib/ldap"
)

// UserAttributes function gets CHESS user attributes
func UserAttributes(user string) (ldap.Entry, error) {
	var attrs ldap.Entry

	// obtain valid token
	httpReadRequest := NewHttpRequest("read", 0)
	httpReadRequest.GetToken()

	// make call to Authz server to obtain user attributes
	rurl := fmt.Sprintf("%s/attrs?user=%s", srvConfig.Config.Services.AuthzURL, user)
	resp, err := httpReadRequest.Get(rurl)
	if err != nil {
		return attrs, err
	}
	// parse data records from meta-data service
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return attrs, err
	}
	err = json.Unmarshal(data, &attrs)
	return attrs, err
}
