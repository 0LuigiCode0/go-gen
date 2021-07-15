package tmp

const HubHelperFuncTmp = `package hub_helper{{$module := .ModuleName}}

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

func InitHelper(H HandlerForHelper) Helper { return &help{HandlerForHelper: H} }

func (h *help) Hash(data string) string {
	hash1 := sha256.New()
	hash2 := md5.New()
	hash1.Write([]byte(data))
	hash2.Write([]byte(data))
	hash1.Write([]byte(string(hash1.Sum(nil)) + string(hash2.Sum(nil))))
	return hex.EncodeToString(hash1.Sum(nil))
}`
