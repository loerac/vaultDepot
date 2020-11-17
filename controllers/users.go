package controllers

import (
    "fmt"
    "net/http"

    "github.com/loerac/vaultDepot/models"
    "github.com/loerac/vaultDepot/views"
)

/**
 * @brief:  Initialize the layout and template for users
 *
 * @return: Users controller
 **/
func NewUsers(userSrv *models.UserService) *Users {
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
func (userRW *Users) NewSignup(writer http.ResponseWriter, request *http.Request) {
    if err := userRW.NewView.Render(writer, nil); err != nil {
        panic(err)
    }
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
func (userRW *Users) CreateSignup(writer http.ResponseWriter, request *http.Request) {
    form := SignupForm{}
    if err := parseForm(request, &form); err != nil {
        panic(err)
    }

    user := models.User {
        FirstName:  form.FirstName,
        LastName:   form.LastName,
        Email:      form.Email,
        Password:   form.Passwd,
    }
    err := userRW.userSrv.Create(&user)
    if err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintln(writer, "Created user:", user)
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
    form := LoginForm{}
    err := parseForm(request, &form)
    if err != nil {
        panic(err)
    }

    user, err := userRW.userSrv.Authenticate(form.Email, form.Passwd)
    switch err {
    case models.ErrNotFound:
        fallthrough
    case models.ErrInvalidPassword:
        fmt.Fprintln(writer, "Email and/or password is incorrect")
    case nil:
        fmt.Fprintln(writer, user)
    default:
        http.Error(writer, err.Error(), http.StatusInternalServerError)
    }
}
