package auth

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
	"gopkg.in/jcmturner/gokrb5.v7/client"
	"gopkg.in/jcmturner/gokrb5.v7/config"
	"gopkg.in/jcmturner/gokrb5.v7/credentials"
)

// Kerberos defines kerberos structure we use
type Kerberos struct {
	User   string
	Ticket []byte
}

// helper function to check user credentials for POST requests
func (k *Kerberos) Credentials() (*credentials.Credentials, error) {
	var msg string
	// user didn't use web interface, we switch to POST form
	fname := fmt.Sprintf("krb-%d", time.Now().UnixNano())
	tmpFile, err := ioutil.TempFile("/tmp", fname)
	if err != nil {
		msg = fmt.Sprintf("Unable to create tempfile: %v", err)
		log.Printf("ERROR: %s", msg)
		return nil, errors.New(msg)
	}
	defer os.Remove(tmpFile.Name())
	_, err = tmpFile.Write(k.Ticket)
	if err != nil {
		msg = "unable to write kerberos ticket"
		log.Printf("ERROR: %s", msg)
		return nil, errors.New(msg)
	}
	err = tmpFile.Close()
	creds, err := kuserFromCache(tmpFile.Name())
	if err != nil {
		msg = "wrong user credentials"
		log.Printf("ERROR: %s", msg)
		return nil, errors.New(msg)
	}
	if creds == nil {
		msg = "unable to obtain user credentials"
		log.Printf("ERROR: %s", msg)
		return nil, errors.New(msg)
	}
	return creds, nil
}

// https://github.com/jcmturner/gokrb5/issues/7
func kuserFromCache(cacheFile string) (*credentials.Credentials, error) {
	cfg, err := config.Load(srvConfig.Config.Kerberos.Krb5Conf)
	ccache, err := credentials.LoadCCache(cacheFile)
	client, err := client.NewClientFromCCache(ccache, cfg)
	err = client.Login()
	if err != nil {
		return nil, err
	}
	return client.Credentials, nil

}
