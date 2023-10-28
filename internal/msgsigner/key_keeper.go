package msgsigner

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

var KeyPath string

type KeyKeeper struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

type KeyConf struct {
	PublicKey  string `mapstructure:"publickey" json:"publickey"`
	PrivateKey string `mapstructure:"privatekey" json:"privatekey"`
}

func init() {
	flag.StringVar(&KeyPath, "keypath", ".", "the key file path")
}

func (k *KeyKeeper) GetPrivateKey() (*rsa.PrivateKey, error) {
	return k.privateKey, nil
}

func (k *KeyKeeper) GetPublicKey() (*rsa.PublicKey, error) {
	return k.publicKey, nil
}

func (k *KeyKeeper) Load() error {
	v := viper.New()
	v.SetConfigFile(KeyPath)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to read config file: %s", err))
	}

	keyConf := KeyConf{}
	err = v.Unmarshal(&keyConf)
	if err != nil {
		panic(fmt.Sprintf("failed to parse key config file: %s", err))
	}

	k.privateKey, err = ParsePrivateKey([]byte(keyConf.PrivateKey))
	if err != nil {
		return err
	}

	k.publicKey, err = ParsePublicKey([]byte(keyConf.PublicKey))
	if err != nil {
		return err
	}

	return nil
}

func ParsePrivateKey(keyBuffer []byte) (*rsa.PrivateKey, error) {
	p, _ := pem.Decode(keyBuffer)
	if p == nil {
		return nil, errors.New("PEM not found")
	}

	return x509.ParsePKCS1PrivateKey(p.Bytes)
}

func ParsePublicKey(keyBuffer []byte) (*rsa.PublicKey, error) {
	p, _ := pem.Decode(keyBuffer)
	if p == nil {
		return nil, errors.New("PEM not found")
	}

	key, err := x509.ParsePKCS1PublicKey(p.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}
