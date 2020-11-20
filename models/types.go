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

    Close() error

    AutoMigrate() error
    DestructiveReset() error
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
    Email       string `gorm:"not null"`
    Username    string
    Website     string `gorm:"not null"`
    Password    string `gorm:"not null"`
}

type VaultService struct {
    db  *gorm.DB
}

var _ UserDB = &userGorm{}
var _ UserService = &userService{}

func (user *User) String() string {
    return fmt.Sprintf("User(Firstname='%s', LastName='%s', Email='%s', Remember='%s')",
        user.FirstName, user.LastName, user.Email, user.RememberHash)
}

func (vault *Vault) String() string {
    return fmt.Sprintf("Vault(Email='%s', Username='%s', Website='%s')",
        vault.Email, vault.Username, vault.Website)
}
