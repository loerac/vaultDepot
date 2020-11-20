package models

import (
    "errors"
)

var (
    /* Return when a resource cannot be found in database */
    ErrNotFound = errors.New("models: Resource not found")

    /* Return when invalid ID is provided */
    ErrIDInvalid = errors.New("models: ID provided was invalid")

    /* Return when password isn't provided */
    ErrPasswordRequired = errors.New("models: Password is required")

    /* Return when password length less than 8 */
    ErrPasswordTooShort = errors.New("models: Password must be at least 8 characters long")

    /* Return when invalid password is used */
    ErrPasswordIncorrect = errors.New("models: Incorrect password")

    /* Return when an email address isn't provided */
    ErrEmailRequired = errors.New("models: Email address is required")

    /* Return when an email address is invalid */
    ErrEmailInvalid = errors.New("models: Email address is not valid")

    /* Return when an provided email is already taken */
    ErrEmailTaken = errors.New("models: Email address is already taken")

    /* Return when a remember token hash isn't provided */
    ErrRememberRequired = errors.New("models: Remember token is required")

    /* Return when remember token isn't at least 32 bytes */
    ErrRememberTooShort = errors.New("models: Remember token must be at least 32 bytes")
)
