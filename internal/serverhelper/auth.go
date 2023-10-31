package serverhelper

import (
	"crypto/md5"
	"encoding/hex"
)

func CalcMsgHash(msg string) string {
	hash := md5.New()
	hash.Write([]byte(msg))
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode
}
