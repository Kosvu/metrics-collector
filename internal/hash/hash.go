package hash

import (
	"crypto/hmac"
	"crypto/sha256"
)

func Sign(data []byte, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	_, err := h.Write(data)

	if err != nil {
		return nil
	}

	res := h.Sum(nil)

	return res
}

func Verify(data, signature, key []byte) bool {
	signature1 := Sign(data, key)

	return hmac.Equal(signature, signature1)
}
