package compat

import (
    "crypto/rand"
    "encoding/base64"
)

const RememberTokenBytes = 32

/**
 * @brief:  Generate n random bytes
 *
 * @param:  n - Amount of bytes to generate
 *
 * @return: [n]bytes on success, else error
 **/
func Bytes(n int) ([]byte, error) {
    bytes := make([]byte, n)
    _, err := rand.Read(bytes)
    if err != nil {
        return nil, err
    }

    return bytes, nil
}

/**
 * @brief:  Generate slice of size n
 *
 * @param:  n - Amount of bytes in slice
 *
 * @return: base64 encoded URL on success, else error
 **/
func String(n int) (string, error) {
    bytes, err := Bytes(n)
    if err != nil {
        return "", err
    }

    return base64.URLEncoding.EncodeToString(bytes), nil
}

/**
 * @brief:  Generate remember token of
 *          predetermined byte size
 *
 * @return: Remember token on success, else error
 **/
func RememberToken() (string, error) {
    return String(RememberTokenBytes)
}
