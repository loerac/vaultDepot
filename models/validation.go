package models

import (
    "strings"

    "github.com/loerac/vaultDepot/compat"
    "github.com/jinzhu/gorm"

    "golang.org/x/crypto/bcrypt"
)

type userValFn func(*User) error
type vaultValFn func(*Vault, User) error

/**
 * @brief:  Iterate over each validation and
 *          normalization function
 *
 * @param:  user - User that is being validated and
 *              normalized
 * @param:  fns - List of functions to run
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

/**
 * @brief:  Iterate over each validation and
 *          normalization function
 *
 * @param:  user - User that is being validated and
 *                 normalized
 * @param:  fns - List of functions to run
 *
 * @return: nil on success, else error
 **/
func runVaultValFns(vault *Vault, user User, fns ...vaultValFn) error {
    for _, fn := range fns {
        if err := fn(vault, user); err != nil {
            return err
        }
    }

    return nil
}

/**
 * @brief:  Create provide user
 *
 * @param:  user - User to create
 *
 * @return: nil on success, else error
 **/
func CreateUser(db *gorm.DB, user *User) error {
    err := runUserValFns(user,
        userPasswordRequired,
        secretKeyRequired,
        passwordMinLength,
        secretKeyMinLength,
        bcryptPassword,
        bcryptSecretKey,
        passwordHashRequired,
        secretKeyHashRequired,
    )
    if err != nil {
        return err
    }

    return db.Create(&user).Error
}

/**
 * @brief:  Create provide user
 *
 * @param:  user - User to create
 *
 * @return: nil on success, else error
 **/
func CreateEntry(db *gorm.DB, vault *Vault, user User) error {
    err := runVaultValFns(vault, user,
        userIDRequired,
        vaultPasswordRequired,
        encryptPassword,
        applicationRequired,
        normalizeApplication,
        normalizeEmail,
        requireEmail,
    )
    if err != nil {
        return err
    }

    return db.Create(&vault).Error
}

/**
 * @brief:  Check if user ID is greater than 0
 *
 * @param:  user - Contains ID
 *
 * @return: nil on success, else ErrIDInvalid
 **/
func userIDRequired(vault *Vault, user User) error {
    if vault.UserID <= 0 {
        return ErrUserIDRequried
    }

    return nil
}

/**
 * @brief:  Checks to see if password is provided
 *
 * @param:  user - Contains password
 *
 * @return: nil on success, else ErrPasswordRequired
 **/
func userPasswordRequired(user *User) error {
    if user.Password == "" {
        return ErrPasswordRequired
    }

    return nil
}

/**
 * @brief:  Checks to see if password is provided
 *
 * @param:  vault - Contains password
 *
 * @return: nil on success, else ErrPasswordRequired
 **/
func vaultPasswordRequired(vault *Vault, user User) error {
    if vault.Password == "" {
        return ErrPasswordRequired
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
func encryptPassword(vault *Vault, user User) error {
    if vault.Password == "" {
        return nil
    }

    aes := compat.NewAES(user.SecretKey)
    passwordCipher, err := aes.Encrypt(vault.Password)
    if err != nil {
        return err
    }

    /* Store ciphered password and forget textbase password */
    vault.PasswordCipher = passwordCipher
    vault.Password = ""

    return nil
}

/**
 * @brief:  Checks to see if secret key is provided
 *
 * @param:  user - Contains secret key
 *
 * @return: nil on success, else ErrSecretKeyRequired
 **/
func secretKeyRequired(user *User) error {
    if user.SecretKey == "" {
        return ErrSecretKeyRequired
    }

    return nil
}

/**
 * @brief:  Check if the password is greater than min length
 *
 * @param:  user - Contains textbased password
 *
 * @return nil on success, else error
 **/
func passwordMinLength(user *User) error {
    if user.Password == "" {
        return nil
    }

    if len(user.Password) < 8 {
        return ErrPasswordTooShort
    }

    return nil
}

/**
 * @brief:  Check if the secret key is greater than min length
 *
 * @param:  user - Contains textbased secret key
 *
 * @return nil on success, else error
 **/
func secretKeyMinLength(user *User) error {
    if user.SecretKey == "" {
        return nil
    }

    if len(user.SecretKey) < 8 {
        return ErrSecretKeyTooShort
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
func bcryptPassword(user *User) error {
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
 * @brief:  Hash a user's secret key, with salt and pepper
 *
 * @param:  user - Contains users secret key
 *
 * @return: nil on success, else error
 **/
func bcryptSecretKey(user *User) error {
    if user.SecretKey == "" {
        return nil
    }

    /* Season textbased secret key with salt and pepper to get hash */
    skBytes := []byte(user.SecretKey + userPwPepper)
    hashedBytes, err := bcrypt.GenerateFromPassword(
        skBytes, bcrypt.DefaultCost,
    )
    if err != nil {
        return err
    }

    /* Store hashed password and forget textbase password */
    user.SecretKeyHash = string(hashedBytes)
    user.SecretKey = ""

    return nil
}

/**
 * @brief:  Checks to see if users password hash has a value
 *
 * @param:  user - Contains password hash
 *
 * @return: nil on success, else error
 **/
func passwordHashRequired(user *User) error {
    if user.PasswordHash == "" {
        return ErrPasswordRequired
    }

    return nil
}

/**
 * @brief:  Checks to see if users secret key hash has a value
 *
 * @param:  user - Contains secret key hash
 *
 * @return: nil on success, else error
 **/
func secretKeyHashRequired(user *User) error {
    if user.SecretKeyHash == "" {
        return ErrSecretKeyRequired
    }

    return nil
}

/**
 * @brief:  Check to see if email is present
 *
 * @param:  vault - Contains email
 *
 * @return: nil on success, else ErrEmailRequired
 **/
func requireEmail(vault *Vault, user User) error {
    if vault.Email == "" {
        return ErrEmailRequired
    }

    return nil
}

/**
 * @brief:  Change email to be lower case and
 *          trim any white spaces in email
 *
 * @param:  vault - Contains email
 *
 * @return: nil
 **/
func normalizeEmail(vault *Vault, user User) error {
    vault.Email = strings.ToLower(vault.Email)
    vault.Email = strings.TrimSpace(vault.Email)
    return nil
}

/**
 * @brief:  Check if application is present
 *
 * @param:  vault - Contains application
 *
 * @return: nil on success, else ErrApplicationRequired
 **/
func applicationRequired(vault *Vault, user User) error {
    if vault.Application == "" {
        return ErrApplicationRequired
    }

    return nil
}

/**
 * @brief:  Change application to be lower case
 *
 * @param:  vault - Contains application
 *
 * @return: nil
 **/
func normalizeApplication(vault *Vault, user User) error {
    vault.Application = strings.ToLower(vault.Application)

    return nil
}
