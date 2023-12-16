package auth

import (
	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/vkuznet/cryptoutils"
)

// LoginForm represents login form
type LoginForm struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// helper function to encrypt login form attributes
func EncryptLoginObject(form LoginForm) (LoginForm, error) {
	encryptedObject, err := cryptoutils.HexEncrypt(
		form.Password, srvConfig.Config.Encryption.Secret, srvConfig.Config.Encryption.Cipher)
	if err != nil {
		return form, err
	} else {
		form.Password = encryptedObject
	}
	return form, nil
}
