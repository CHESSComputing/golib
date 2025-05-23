package utils

import (
	"encoding/json"
	"os"

	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/vkuznet/cryptoutils"
)

// TrustedClient represents trusted client
type TrustedClient struct {
	IPs  []string           `json:"ip_addresses"`
	User string             `json:"user"`
	MACs []MacAddressRecord `json:"mac_addresses"`
}

// NewTrustedClient provides pointer to trusted client initialized with appropriate fields
func NewTrustedClient() *TrustedClient {
	t := &TrustedClient{
		IPs:  IpAddr(),
		MACs: MacAddr(),
		User: os.Getenv("USER"),
	}
	return t
}

// Encrypt encrypt trusted client information
func (t *TrustedClient) Encrypt(salt string) ([]byte, error) {
	var edata []byte
	data, err := json.Marshal(t)
	if err != nil {
		return edata, err
	}
	cipher := srvConfig.Config.Encryption.Cipher
	return cryptoutils.Encrypt(data, salt, cipher)
}

// Decrypt decrypt trusted client information
func (t *TrustedClient) Decrypt(edata []byte, salt string) error {
	var tdata TrustedClient
	cipher := srvConfig.Config.Encryption.Cipher
	data, err := cryptoutils.Decrypt(edata, salt, cipher)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &tdata)
	t.User = tdata.User
	t.IPs = tdata.IPs
	t.MACs = tdata.MACs
	return err
}
