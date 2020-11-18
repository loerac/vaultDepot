package compat

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "hash"
)

type HMAC struct {
    hmac    hash.Hash
}

/**
 * @brief:  Create new HMAC
 *
 * @param:  key - Used to sign data
 *
 * @return: HMAC on success
 **/
func NewHMAC(key string) HMAC {
    hmacObj := hmac.New(sha256.New, []byte(key))
    return HMAC {
        hmac: hmacObj,
    }
}

/**
 * @brief:  Hash input with secrect key
 *
 * @param:  input - String to hash using HMAC
 *
 * @return: base64 Encoded URL as string
 **/
func (hmacObj HMAC) Hash(input string) string {
    hmacObj.hmac.Reset()
    hmacObj.hmac.Write([]byte(input))
    bytes := hmacObj.hmac.Sum(nil)

    return base64.URLEncoding.EncodeToString(bytes)
}
