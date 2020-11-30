package compat

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/md5"
    "crypto/rand"
    "encoding/hex"
    "io"
)

type AES struct {
    aes     string
}

/**
 * @brief:  Create a 32-bit length key
 *
 * @param:  key - User key to salt user data
 *
 * @return: Key as a 32-bit length md5sum
 **/
func NewAES(key string) AES {
    hasher := md5.New()
    hasher.Write([]byte(key))

    return AES {
        aes:    hex.EncodeToString(hasher.Sum(nil)),
    }
}

/**
 * @brief:  Encrypt the given input with the users key
 *
 * @param:  data - Input that is being encrypted
 *
 * @return: Data encrypted
 **/
func (aesObj AES) Encrypt(data string) ([]byte, error) {
    block, _ := aes.NewCipher([]byte(aesObj.aes))
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    return gcm.Seal(nonce, nonce, []byte(data), nil), nil
}

/**
 * @brief:  Decrypt the given input with the users key
 *
 * @param:  data - Input that is being decrypted
 *
 * @return: Data decrypted
 **/
func (aesObj AES) Decrypt(data []byte) (string, error) {
    block, err := aes.NewCipher([]byte(aesObj.aes))
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonceSize := gcm.NonceSize()
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}
