package models

import (
    "fmt"

    "github.com/jinzhu/gorm"
)

const (
    userPwPepper = "secret-random-string"
)

type User struct {
    gorm.Model
    Username    string `gorm:"not null;unique_index"`
    Password    string `gorm:"-"`
    PasswordHash string `gorm:"not null"`
    SecretKey   string `gorm:"-"`
    SecretKeyHash string `gorm:"not null"`
}

type Vault struct {
    gorm.Model
    UserID      uint `gorm:"not_null;index"`
    Email       string `gorm:"not null"`
    Username    string
    Application string `gorm:"not null"`
    Password    string `gorm:"-"`
    PasswordCipher []byte `gorm:"not null"`
}

func (user User) String() string {
    return fmt.Sprintf("User(Username='%s', Password='%s', SecretKey='%s')",
        user.Username, user.PasswordHash, user.SecretKeyHash)
}

func (vault Vault) String() string {
    return fmt.Sprintf("Vault(Email='%s', Username='%s', Application='%s')",
        vault.Email, vault.Username, vault.Application)
}
