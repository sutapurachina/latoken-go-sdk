package latoken_go_sdk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
)

func GetSignature(secretKey string, msg []byte) []byte {
	signer := hmac.New(sha256.New, bytes.NewBufferString(secretKey).Bytes())
	signer.Write(msg)
	return signer.Sum(nil)
}
