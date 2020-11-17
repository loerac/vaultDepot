package models

import (
    "fmt"

    "github.com/loerac/vaultDepot/compat"

    "github.com/jinzhu/gorm"
)

const (
    userPwPepper = "secret-random-string"
    hmacSecretKey = "secret-hmac-key"
)

type User struct {
    gorm.Model
    FirstName   string
    LastName    string
    Email       string `gorm:"not null;unique_index"`
    Password    string `gorm:"-"`
    PasswordHash string `gorm:"not null"`
    Remember    string `gorm:"-"`
    RememberHash string `gorm:"not null;unique_index"`
}

type UserService struct {
    db      *gorm.DB
    hmac    compat.HMAC
}

type Vault struct {
    gorm.Model
    Email       string `gorm:"not null"`
    Username    string
    Website     string `gorm:"not null"`
    Password    string `gorm:"not null"`
}

type VaultService struct {
    db  *gorm.DB
}

func (user *User) String() string {
    return fmt.Sprintf("User(Firstname='%s', LastName='%s', Email='%s', Remember='%s')",
        user.FirstName, user.LastName, user.Email, user.RememberHash)
}

func (vault *Vault) String() string {
    return fmt.Sprintf("Vault(Email='%s', Username='%s', Website='%s')",
        vault.Email, vault.Username, vault.Website)
}
