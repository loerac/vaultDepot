package controllers

import (
    "fmt"
    "net/http"

    "github.com/loerac/vaultDepot/compat"
    "github.com/loerac/vaultDepot/models"
    "github.com/loerac/vaultDepot/views"
)

/**
 * @brief:  Initialize the layout and template for users
 *
 * @return: Users controller
 **/
func NewUsers(userSrv models.UserService) *Users {
    return &Users {
        NewView: views.NewView("bootstrap", "users/new"),
        LoginView: views.NewView("bootstrap", "users/login"),
        userSrv: userSrv,
    }
}

/* ==============================*/
/*     METHODS FOR RESTful API   */
/* ==============================*/

/**
 * @brief:  Render form where a user can
 *          create a new user account
 *
 * @param:  writer - Render the template
 * @param:  request - Input data from user
 *
 * @action: GET /signup
 **/
func (userRW *Users) New(writer http.ResponseWriter, request *http.Request) {
    userRW.NewView.Render(writer, request, nil)
}

/**
 * @brief:  Process the signup form when a user
 *          tries to create a new user account
 *
 * @param:  writer - Render the html template
 * @param:  request - Input data from user
 *
 * @action: POST /signup
 **/
func (userRW *Users) Create(writer http.ResponseWriter, request *http.Request) {
    viewData := views.Data{}
    form := SignupForm{}
    if err := parseForm(request, &form); err != nil {
        viewData.SetAlert(err)
        userRW.NewView.Render(writer, request, viewData)
        return
    }

    user := models.User {
        FirstName:  form.FirstName,
        LastName:   form.LastName,
        Email:      form.Email,
        Password:   form.Passwd,
        SecretKey:  form.SecretKey,
    }
    err := userRW.userSrv.Create(&user)
    if err != nil {
        viewData.SetAlert(err)
        userRW.NewView.Render(writer, request, viewData)
        return
    }

    err = userRW.signIn(writer, &user)
    if err != nil {
        http.Redirect(writer, request, "/login", http.StatusFound)
        return
    }

    http.Redirect(writer, request, "/vault/new", http.StatusFound)
}

/**
 * @brief:  Process the login form when a user tries to log in
 *
 * @param:  writer - Render the html template
 * @param:  request - Input data from user
 *
 * @action: POST /login
 **/
func (userRW *Users) Login(writer http.ResponseWriter, request *http.Request) {
    viewData := views.Data{}
    form := LoginForm{}
    err := parseForm(request, &form)
    if err != nil {
        viewData.SetAlert(err)
        userRW.LoginView.Render(writer, request, viewData)
        return
    }

    user, err := userRW.userSrv.Authenticate(form.Email, form.Passwd)
    if err != nil {
        switch err {
        case models.ErrNotFound:
            fallthrough
        case models.ErrPasswordIncorrect:
            viewData.SetupAlert(views.AlertError, "Email and/or password is incorrect")
        default:
            viewData.SetAlert(err)
        }
        userRW.LoginView.Render(writer, request, viewData)

        return
    }

    err = userRW.signIn(writer, user)
    if err != nil {
        viewData.SetAlert(err)
        userRW.LoginView.Render(writer, request, viewData)
        return
    }

    http.Redirect(writer, request, "/vault/account", http.StatusFound)
}

/* ==============================*/
/*          MISC METHODS         */
/* ==============================*/

/**
 * @brief:  Sign the given user in via cookies
 *
 * @param:  writer - Set user's cookie
 * @param:  user - Update remember token
 *
 * @return: nil on success, else error
 **/
func (userRW *Users) signIn(writer http.ResponseWriter, user *models.User) error {
    if user.Remember == "" {
        token, err := compat.RememberToken()
        if err != nil {
            return err
        }

        user.Remember = token
        err = userRW.userSrv.Update(user)
        if err != nil {
            return err
        }
    }

    cookie := http.Cookie {
        Name:       "remember_token",
        Value:      user.Remember,
        HttpOnly:   true,
    }
    http.SetCookie(writer, &cookie)

    return nil
}

/* Dev testing: display cookies set on the current user */
func (userRW *Users) CookiesTest(writer http.ResponseWriter, request *http.Request) {
    cookie, err := request.Cookie("remember_token")
    if err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
        return
    }

    user, err := userRW.userSrv.ByRemember(cookie.Value)
    if err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintln(writer, user)
}
