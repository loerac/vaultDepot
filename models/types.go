package models

import (
    "fmt"
    "regexp"

    "github.com/loerac/vaultDepot/compat"

    "github.com/jinzhu/gorm"
)

const (
    userPwPepper = "secret-random-string"
    hmacSecretKey = "secret-hmac-key"
)

type UserService interface {
    Authenticate(email, password string) (*User, error)
    UserDB
}

type UserDB interface {
    ByID(id uint) (*User, error)
    ByEmail(email string) (*User, error)
    ByRemember(token string) (*User, error)

    Create(user *User) error
    Update(user *User) error
    Delete(id uint) error
}

type User struct {
    gorm.Model
    FirstName   string
    LastName    string
    Email       string `gorm:"not null;unique_index"`
    Password    string `gorm:"-"`
    PasswordHash string `gorm:"not null"`
    Remember    string `gorm:"-"`
    RememberHash string `gorm:"not null;unique_index"`
    SecretKey   string `gorm:"not null"`
}

type userService struct {
    UserDB
}

type userGorm struct {
    db      *gorm.DB
}

type userValidator struct {
    UserDB
    hmac    compat.HMAC
    emailRegex *regexp.Regexp
}

type Vault struct {
    gorm.Model
    UserID      uint   `gorm:"not_null;index"`
    SecretKey   string `gorm:"-"`
    Email       string `gorm:"not null"`
    Username    string
    Application string `gorm:"not null"`
    Password    string `gorm:"-"`
    PasswordHash []byte `gorm:"not null"`
}

type VaultDB interface {
    ByID(id uint, key string) (*[]Vault, error)
    ByUserID(userID uint) ([]Vault, error)
    Create(vault *Vault) error
}

type VaultService interface {
    VaultDB
}

type vaultGorm struct {
    db          *gorm.DB
}

type vaultService struct {
    VaultDB
}

type vaultValidator struct {
    VaultDB
}

type Services struct {
    db          *gorm.DB
    User        UserService
    Vault       VaultService
}

var _ UserDB = &userGorm{}
var _ VaultDB = &vaultGorm{}
var _ UserService = &userService{}

func (user *User) String() string {
    return fmt.Sprintf("User(Firstname='%s', LastName='%s', Email='%s', Remember='%s', SecretKey='%s')",
        user.FirstName, user.LastName, user.Email, user.RememberHash, user.SecretKey)
}

func (vault *Vault) String() string {
    return fmt.Sprintf("Vault(Email='%s', Username='%s', Application='%s', SecretKey='%s')",
        vault.Email, vault.Username, vault.Application, vault.SecretKey)
}
