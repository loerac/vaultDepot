package models

import (
    "fmt"

    "github.com/loerac/vaultDepot/compat"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    "golang.org/x/crypto/bcrypt"
)

/**
 * @brief:  Initialize user database with GORM, and
 *          initialize user validation
 *
 * @param:  connInfo - Information of database
 *
 * @return: UserService on success, else error
 **/
func NewUserService(connInfo string) (UserService, error) {
    db, err := newUserGorm(connInfo)
    if err != nil {
        return nil, err
    }

    hmac := compat.NewHMAC(hmacSecretKey)
    userValid := newUserValidator(db, hmac)

    return &userService{
        UserDB: userValid,
    }, nil
}

/**
 * @brief:  Connecting to user database with GORM
 *
 * @param:  connInfo - Information of database
 *
 * @return: userGorm on success, else error
 **/
func newUserGorm(connInfo string) (*userGorm, error) {
    db, err := gorm.Open("postgres", connInfo)
    if err != nil {
        return nil, err
    }
    db.LogMode(true)

    return &userGorm{
        db:     db,
    }, nil
}

/**
 * @brief:  Close the database connection
 *
 * @return: nil on success, else error
 **/
func (usergorm *userGorm) Close() error {
    return usergorm.db.Close()
}

/* ==============================*/
/*        METHODS FOR CRUD       */
/* ==============================*/

/**
 * @brief:  Create provided user
 *
 * @param:  user - User information as struct
 *
 * @return: nil on success, else error
 **/
func (usergorm *userGorm) Create(user *User) error {
    return usergorm.db.Create(user).Error
}

/**
 * @brief:  Update provided user with data
 *
 * @param:  user - User information
 *
 * @return: nil on success, else error
 **/
func (usergorm *userGorm) Update(user *User) error {
    return usergorm.db.Save(user).Error
}

/**
 * @brief:  Delete provided user with ID
 *
 * @param:  id - User to be deleted
 *
 * @return: nil on success
 *          ErrIDInvalid if ID is invalid
 *          Else error
 **/
func (usergorm *userGorm) Delete(id uint) error {
    user := User{
        Model: gorm.Model {
            ID: id,
        },
    }
    return usergorm.db.Delete(&user).Error
}

/* ==============================*/
/*   METHODS TO SEARCH FOR USER  */
/* ==============================*/

/**
 * @brief:  Look up a user with provided ID.
 *
 * @param:  id  - ID of the user
 *
 * @return: If user is found, return nil
 *          If user not found, return ErrNotFound
 *          Else, return error
 **/
func (usergorm *userGorm) ByID(id uint) (*User, error) {
    var user User
    db := usergorm.db.Where("id = ?", id)
    err := first(db, &user)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

/**
 * @brief:  Look up a user with provided email addr
 *
 * @param:  email - Email address of user
 *
 * @return: If user is found, return nil
 *          If user not found, return ErrNotFound
 *          Else, return error
 **/
func (usergorm *userGorm) ByEmail(email string) (*User, error) {
    var user User
    db := usergorm.db.Where("email = ?", email)
    err := first(db, &user)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

/**
 * @brief:  Look up a user with given remember token
 *
 * @param:  rememberHash - Token of user already hashed
 *
 * @return: If user is found, return nil
 *          If user not found, return ErrNotFound
 *          Else, return error
 **/
func (usergorm *userGorm) ByRemember(rememberHash string) (*User, error) {
    var user User
    db := usergorm.db.Where("remember_hash = ?", rememberHash)
    err := first(db, &user)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

/**
 * @brief:  Authenticate a user with provided
 *          email addr and password
 *
 * @param:  email - Users email address in DB
 * @param:  password - Users password for account
 *
 * @return: If email addr is invalid, return ErrNotFound
 *          If password is invalid, return ErrPasswordIncorrect
 *          If both are vaild, return user
 *          Else, error
 **/
func (usergorm *userService) Authenticate(email, password string) (*User, error) {
    foundUser, err := usergorm.ByEmail(email)
    if err != nil {
        return nil, err
    }

    err = bcrypt.CompareHashAndPassword(
        []byte(foundUser.PasswordHash),
        []byte(password + userPwPepper),
    )
    switch err {
    case nil:
        return foundUser, nil
    case bcrypt.ErrMismatchedHashAndPassword:
        return nil, ErrPasswordIncorrect
    default:
        return nil, err
    }
}

/* ==============================*/
/*      METHODS FOR DATABASE     */
/*           MIGRATION           */
/* ==============================*/

/**
 * @brief:  Attempt to automatically migrate user table
 *
 * @return: nil on success, else error
 **/
func (usergorm *userGorm) AutoMigrate() error {
    err := usergorm.db.AutoMigrate(&User{}).Error
    if err != nil {
        fmt.Println("models: Error on migrating User table: ", err)
    }

    return err
}

/**
 * @brief:  Drops the user table and rebuilds it
 **/
func (usergorm *userGorm) DestructiveReset() error {
    err := usergorm.db.DropTableIfExists(&User{}).Error
    if err != nil {
        fmt.Println("models: Error on dropping User table: ", err)
        return err
    }

    return usergorm.AutoMigrate()
}
