package xcache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func GenerateKey(prefix string, id string) string {
	return fmt.Sprintf("%s%s", prefix, id)
}

func GenerateKeyWithParams(prefix string, params ...interface{}) string {
	key := prefix
	for _, param := range params {
		key = fmt.Sprintf("%s:%v", key, param)
	}
	return key
}

func HashKey(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func BuildPattern(prefix string) string {
	return fmt.Sprintf("%s*", prefix)
}
