package models

import (
	"bufio"
	"fmt"
	"os"
	"strings"

    "github.com/loerac/vaultDepot/compat"
    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"

    "golang.org/x/crypto/bcrypt"
)

/**
 * @brief:  Create a new user
 *
 * @param:   db - Pointer to database
 *
 * @return: On success, a new user
 *          Else, an error
 **/
func Signup(db *gorm.DB) (User, error) {
    username := ""
    for username == "" {
        fmt.Print("Enter username: ")
        fmt.Scanln(&username)
    }

    password := HiddenInput("password")
    secret_key := HiddenInput("secret key")

    new_user := User {
        Username: username,
        Password: password,
        SecretKey: secret_key,
    }

    if err := CreateUser(db, &new_user); err != nil {
        return User{}, err
    }

    return new_user, nil
}

/**
 * @brief:  Log user in
 *
 * @param:   db - Pointer to database
 *
 * @return: On success, user
 *          Else, an error
 **/
func Login(userdb *gorm.DB) (User, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')
    username = strings.TrimSpace(username)

    password, err := UserInput("Enter Password")
	if err != nil {
        return User{}, err
	}

    secret_key, err := UserInput("Enter secret key")
	if err != nil {
        return User{}, err
	}
    fmt.Println()

    login_user := User {
        Username: username,
        Password: password,
        SecretKey: secret_key,
    }

    user, err := Authenticate(userdb, login_user)
    compat.CheckError(err)

	return *user, nil
}

/**
 * @brief:  Look up a user with provided email addr
 *
 * @param:  db - Pointer to database
 * @param:  username - User to look up
 *
 * @return: If user is found, return user
 *          If user not found, return ErrNotFound
 *          Else, return error
 **/
func ByUsername(userdb *gorm.DB, username string) (*User, error) {
    var user User
    db := userdb.Where("username = ?", username)
    err := first(db, &user)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

/**
 * @brief:  Authenticate a user with provided
 *          username and password
 *
 * @param:  userdb - db pointer to database
 * @param:  auth_user - User to authenticate
 *
 * @return: If username is invalid, return ErrNotFound
 *          If password is invalid, return ErrPasswordIncorrect
 *          If both are vaild, return user
 *          Else, error
 **/
func Authenticate(userdb *gorm.DB, auth_user User) (*User, error) {
    foundUser, err := ByUsername(userdb, auth_user.Username)
    if err != nil {
        return nil, err
    }

    password_err := bcrypt.CompareHashAndPassword(
        []byte(foundUser.PasswordHash),
        []byte(auth_user.Password + userPwPepper),
    )
    secret_key_err := bcrypt.CompareHashAndPassword(
        []byte(foundUser.SecretKeyHash),
        []byte(auth_user.SecretKey + userPwPepper),
    )
    if password_err == bcrypt.ErrMismatchedHashAndPassword {
        return nil, ErrPasswordIncorrect
    } else if secret_key_err == bcrypt.ErrMismatchedHashAndPassword {
        return nil, ErrSecretKeyIncorrect
    } else if password_err == nil &&
              secret_key_err == nil {
        return foundUser, nil
    }

    return nil, err
}
