package middleware

import (
    "net/http"

    "github.com/loerac/vaultDepot/context"
    "github.com/loerac/vaultDepot/models"
)

type User struct {
    models.UserService
}

type RequireUser struct {}

/**
 * @brief:  Check to see if a user is logged in.
 *          If logged in, call next.
 *          Else, redirect to sign in
 *
 * @param:  writer - Render the html template
 * @param:  request - Input data from user
 *
 * @return: HandlerFunc to next page or sign in
 **/
func (reqUser *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
        user := context.User(request.Context())
        if user == nil {
            http.Redirect(writer, request, "/login", http.StatusFound)
            return
        }

        next(writer, request)
    })
}

/**
 * @brief:  Get the remember token cookie and look up the user
 *
 * @param:  writer - Render the html template
 * @param:  request - Input data from user
 *
 * @return: HandlerFunc to next page or sign in
 **/
func (user *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
        cookie, err := request.Cookie("remember_token")
        if err != nil {
            next(writer, request)
            return
        }

        user, err := user.UserService.ByRemember(cookie.Value)
        if err != nil {
            next(writer, request)
            return
        }

        ctx := request.Context()
        ctx = context.WithUser(ctx, user)
        request = request.WithContext(ctx)
        next(writer, request)
    })
}

func (reqUser *RequireUser) Apply(next http.Handler) http.HandlerFunc {
    return reqUser.ApplyFn(next.ServeHTTP)
}

func (user *User) Apply(next http.Handler) http.HandlerFunc {
    return user.ApplyFn(next.ServeHTTP)
}
