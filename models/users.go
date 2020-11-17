package models

import (
    "fmt"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    "golang.org/x/crypto/bcrypt"
)

/**
 * @brief:  Connecting to user database with GORM
 *
 * @param:  connInfo - Information of database
 *
 * @return: UserService on success, else error
 **/
func NewUserService(connInfo string) (*UserService, error) {
    db, err := gorm.Open("postgres", connInfo)
    if err != nil {
        return nil, err
    }
    db.LogMode(true)

    return &UserService{
        db: db,
    }, nil
}

/**
 * @brief:  Close the database connection
 *
 * @return: nil on success, else error
 **/
func (userSrv *UserService) Close() error {
    return userSrv.db.Close()
}

/* ==============================*/
/*        METHODS FOR CRUD       */
/* ==============================*/

/**
 * @brief:  Create provided user and backfill data
 *
 * @param:  user - User information as struct
 *
 * @return: nil on success, else error
 **/
func (userSrv *UserService) Create(user *User) error {
    /* Season textbased password with salt and pepper to get hash */
    pwBytes := []byte(user.Password + userPwPepper)
    hashedBytes, err := bcrypt.GenerateFromPassword(
        pwBytes, bcrypt.DefaultCost,
    )
    if err != nil {
        return err
    }

    /* Store hashed password and forget textbase password */
    user.PasswordHash = string(hashedBytes)
    user.Password = ""

    return userSrv.db.Create(user).Error
}

/**
 * @brief:  Update provided user with data
 *
 * @param:  user - User information
 *
 * @return: nil on success, else error
 **/
func (userSrv *UserService) Update(user *User) error {
    return userSrv.db.Save(user).Error
}

/**
 * @brief:  Delete provided user with ID
 *
 * @param:  id - User to be deleted
 *
 * @return: nil on success
 *          ErrInvalidID if ID is invalid
 *          Else error
 **/
func (userSrv *UserService) Delete(id uint) error {
    if id == 0 {
        return ErrInvalidID
    }

    user := User{
        Model: gorm.Model {
            ID: id,
        },
    }
    return userSrv.db.Delete(&user).Error
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
func (userSrv *UserService) ByID(id uint) (*User, error) {
    var user User
    db := userSrv.db.Where("id = ?", id)
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
func (userSrv *UserService) ByEmail(email string) (*User, error) {
    var user User
    db := userSrv.db.Where("email = ?", email)
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
 *          If password is invalid, return ErrInvalidPassword
 *          If both are vaild, return user
 *          Else, error
 **/
func (userSrv *UserService) Authenticate(email, password string) (*User, error) {
    foundUser, err := userSrv.ByEmail(email)
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
        return nil, ErrInvalidPassword
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
func (userSrv *UserService) AutoMigrate() error {
    err := userSrv.db.AutoMigrate(&User{}).Error
    if err != nil {
        fmt.Println("models: Error on migrating User table: ", err)
    }

    return err
}

/**
 * @brief:  Drops the user table and rebuilds it
 **/
func (userSrv *UserService) DestructiveReset() error {
    err := userSrv.db.DropTableIfExists(&User{}).Error
    if err != nil {
        fmt.Println("models: Error on dropping User table: ", err)
        return err
    }

    return userSrv.AutoMigrate()
}
