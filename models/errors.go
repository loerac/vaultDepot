package models

import (
    "strings"
)

type modelError string
type privateError string

var (
    /* Return when a resource cannot be found in database */
    ErrNotFound modelError = "models: Resource not found"

    /* Return when invalid ID is provided */
    ErrIDInvalid modelError = "models: ID provided was invalid"

    /* Return when password isn't provided */
    ErrPasswordRequired modelError = "models: Password is required"

    /* Return when password length less than 8 */
    ErrPasswordTooShort modelError = "models: Password must be at least 8 characters long"

    /* Return when invalid password is used */
    ErrPasswordIncorrect modelError = "models: Incorrect password"

    /* Return when an email address isn't provided */
    ErrEmailRequired modelError = "models: Email address is required"

    /* Return when an email address is invalid */
    ErrEmailInvalid modelError = "models: Email address is not valid"

    /* Return when an provided email is already taken */
    ErrEmailTaken modelError = "models: Email address is already taken"

    /* Return when a remember token hash isn't provided */
    ErrRememberRequired privateError = "models: Remember token is required"

    /* Return when remember token isn't at least 32 bytes */
    ErrRememberTooShort privateError = "models: Remember token must be at least 32 bytes"
)

func (err modelError) Error() string {
    return string(err)
}

func (err modelError) Public() string {
    str := strings.Replace(string(err), "models: ", "", 1)
    return str
}

func (err privateError) Error() string {
    return string(err)
}
