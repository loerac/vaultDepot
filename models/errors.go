package models

import (
    "errors"
)

var (
    /* Return when a resource cannot be found in database */
    ErrNotFound = errors.New("models: resource not found")

    /* Return when invalid ID is provided */
    ErrInvalidID = errors.New("models: ID provided was invalid")

    /* Return when invalid password is used */
    ErrInvalidPassword = errors.New("models: incorrect password")
)
