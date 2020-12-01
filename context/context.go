package context

import (
    "context"

    "github.com/loerac/vaultDepot/models"
)

type privateKey string

const (
    userKey privateKey = "user"
)

/**
 * @brief:  Assign a value to the "user" key for a user
 *
 * @param:  ctx - http.Request context for current handler
 * @param:  user - Value for the "user" key
 *
 * @return: New context with the user set as a value
 **/
func WithUser(ctx context.Context, user *models.User) context.Context {
    return context.WithValue(ctx, userKey, user)
}

/**
 * @brief:  Retrieve currect user logged in with context
 *
 * @param:  ctx - http.Request context for current handler
 *
 * @return: Current user logged in
 **/
func User(ctx context.Context) *models.User {
    if temp := ctx.Value(userKey); temp != nil {
        if user, ok := temp.(*models.User); ok {
            return user
        }
    }

    return nil
}
