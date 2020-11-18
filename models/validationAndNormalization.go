package models

import (
    "regexp"
    "strings"

    "github.com/loerac/vaultDepot/compat"

    "golang.org/x/crypto/bcrypt"
)

type userValFn func(*User) error

/**
 * @brief:  Initialize user validation
 *
 * @param:  userDB - User database connection
 * @param:  hmac - user HMAC
 *
 * @return: User Validator
 **/
func newUserValidator(userDB UserDB, hmac compat.HMAC) *userValidator {
    return &userValidator {
        UserDB: userDB,
        hmac:   hmac,
        emailRegex: regexp.MustCompile(
            `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
    }
}

/**
 * @brief:  Iterate over each validation and
 *          normalization function
 *
 * @param:  user - User that is being validated and
 *              normalized
 *
 * @return: nil on success, else error
 **/
func runUserValFns(user *User, fns ...userValFn) error {
    for _, fn := range fns {
        if err := fn(user); err != nil {
            return err
        }
    }

    return nil
}

/* ==============================*/
/*      METHODS TO VALIDATION    */
/*        AND NORMALIZATION      */
/* ==============================*/

/**
 * @brief:  Create provide user
 *
 * @param:  user - User to create
 *
 * @return: nil on success, else error
 **/
func (userValid *userValidator) Create(user *User) error {
    err := runUserValFns(user,
        userValid.passwordRequired,
        userValid.passwordMinLength,
        userValid.bcryptPassword,
        userValid.passwordHashRequired,
        userValid.setRememberIfUnset,
        userValid.rememberMinBytes,
        userValid.hmacRemember,
        userValid.rememberHashRequired,
        userValid.normalizeEmail,
        userValid.requireEmail,
        userValid.emailFormat,
        userValid.emailIsAvail,
    )
    if err != nil {
        return err
    }

    return userValid.UserDB.Create(user)
}

/**
 * @brief:  Update provided user with data
 *
 * @param:  user - User information
 *
 * @return: nil on success, else error
 **/
func (userValid *userValidator) Update(user *User) error {
    err := runUserValFns(user,
        userValid.passwordMinLength,
        userValid.bcryptPassword,
        userValid.passwordHashRequired,
        userValid.rememberMinBytes,
        userValid.hmacRemember,
        userValid.rememberHashRequired,
        userValid.normalizeEmail,
        userValid.requireEmail,
        userValid.emailFormat,
        userValid.emailIsAvail,
    )
    if err != nil {
        return err
    }

    return userValid.UserDB.Update(user)
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
func (userValid *userValidator) Delete(id uint) error {
    var user User
    user.ID = id

    err := runUserValFns(&user,
        userValid.validId,
    )
    if err != nil {
        return err
    }

    return userValid.UserDB.Delete(id)
}

/**
 * @brief:  Hash the remember token
 *
 * @param:  token - Token of user
 *
 * @return: User on success, else error
 **/
func (userValid *userValidator) ByRemember(token string) (*User, error) {
    user := User {
        Remember: token,
    }

    err := runUserValFns(&user,
        userValid.hmacRemember,
    )
    if err != nil {
        return nil, err
    }

    return userValid.UserDB.ByRemember(user.RememberHash)
}

/**
 * @brief:  Normalize email address
 *
 * @param:  email - Email address that is going to be normalized
 *
 * @return: User on success, else error
 **/
func (userValid *userValidator) ByEmail(email string) (*User, error) {
    user := User {
        Email:  email,
    }

    err := runUserValFns(&user, userValid.normalizeEmail)
    if err != nil {
        return nil, err
    }

    return userValid.UserDB.ByEmail(user.Email)
}

/**
 * @brief:  Checks to see if password is provided
 *
 * @param:  user - Contains password
 *
 * @return: nil on success, else ErrPasswordRequired
 **/
func (userValid *userValidator) passwordRequired(user *User) error {
    if user.Password == "" {
        return ErrPasswordRequired
    }

    return nil
}

/**
 * @brief:  Checks to see if users password hash has a value
 *
 * @param:  user - Contains password hash
 *
 * @return: nil on success, else error
 **/
func (userValid *userValidator) passwordHashRequired(user *User) error {
    if user.PasswordHash == "" {
        return ErrPasswordRequired
    }

    return nil
}

/**
 * @brief:  Check if the password is greaterc than min length
 *
 * @param:  user - Contains textbased password
 *
 * @return nil on success, else error
 **/
func (userValid *userValidator) passwordMinLength(user *User) error {
    if user.Password == "" {
        return nil
    }

    if len(user.Password) < 8 {
        return ErrPasswordTooShort
    }

    return nil
}

/**
 * @brief:  Hash a user's password, with salt and pepper
 *
 * @param:  user - Contains users password
 *
 * @return: nil on success, else error
 **/
func (userValid *userValidator) bcryptPassword(user *User) error {
    if user.Password == "" {
        return nil
    }

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

    return nil
}

/**
 * @brief:  Check the number of bytes is greater than 32
 *
 * @param:  user - Contains remember token
 *
 * @return: nil on success, else ErrRememberTooShort
 **/
func (userValid *userValidator) rememberMinBytes(user *User) error {
    if user.Remember == "" {
        return nil
    }

    n, err := compat.NBtyes(user.Remember)
    if err != nil {
        return err
    }

    if n < 32 {
        return ErrRememberTooShort
    }

    return nil
}

/**
 * @brief:  Check remember token hash is set
 *
 * @param:  user - Contains remember token hash
 *
 * @return: nil on success, else ErrRememberRequired
 **/
func (userValid *userValidator) rememberHashRequired(user *User) error {
    if user.RememberHash == "" {
        return ErrRememberRequired
    }

    return nil
}

/**
 * @brief:  Hash the remember token if isn't present
 *
 * @param:  user - Contains token
 *
 * @return: nil
 **/
func (userValid *userValidator) hmacRemember(user *User) error {
    if user.Remember == "" {
        return nil
    }

    user.RememberHash = userValid.hmac.Hash(user.Remember)
    return nil
}

/**
 * @brief:  Set a remember token if one is present
 *
 * @user:   user - Contains the remember token
 *
 * @return: nil on success, else error
 **/
func (userValid *userValidator) setRememberIfUnset(user *User) error {
    if user.Remember != "" {
        return nil
    }

    token, err := compat.RememberToken()
    if err != nil {
        return nil
    }

    user.Remember = token
    return nil
}

/**
 * @brief:  Check if user ID is greater than 0
 *
 * @param:  user - Contains ID
 *
 * @return: nil on success, else ErrIDInvalid
 **/
func (userValid *userValidator) validId(user *User) error {
    if user.ID <= 0 {
        return ErrIDInvalid
    }

    return nil
}

/**
 * @brief:  Change email to be lower case and
 *          trim any white spaces in email
 *
 * @param:  user - Contains email
 *
 * @return: nil
 **/
func (userValid *userValidator) normalizeEmail(user *User) error {
    user.Email = strings.ToLower(user.Email)
    user.Email = strings.TrimSpace(user.Email)
    return nil
}

/**
 * @brief:  Check to see if email is present
 *
 * @param:  user - Contains email
 *
 * @return: nil on success, else ErrEmailRequired
 **/
func (userValid *userValidator) requireEmail(user *User) error {
    if user.Email == "" {
        return ErrEmailRequired
    }

    return nil
}

/**
 * @brief:  Check if the user email matches regular expression
 *
 * @param:  user - Contains email
 *
 * @return: nil on success, else ErrEmailInvalid
 **/
func (userValid *userValidator) emailFormat(user *User) error {
    if user.Email == "" {
        return nil
    }

    if !userValid.emailRegex.MatchString(user.Email) {
        return ErrEmailInvalid
    }

    return nil
}

/**
 * @brief:  check if email is taken
 *
 * @param:  user - Contains email
 *
 * @return: nil on success, else error
 **/
func (userValid *userValidator) emailIsAvail(user *User) error {
    foundUser, err := userValid.ByEmail(user.Email)
    if err == ErrNotFound {
        /* Email is available */
        return nil
    }

    if err != nil {
        return err
    }

    /* Check if found user's ID match user provided */
    if user.ID != foundUser.ID {
        return ErrEmailTaken
    }

    return nil
}
