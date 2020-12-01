package controllers

import (
    "github.com/loerac/vaultDepot/models"
    "github.com/loerac/vaultDepot/views"

    "github.com/gorilla/mux"
)

type Users struct {
    NewView     *views.View
    LoginView   *views.View
    userSrv     models.UserService
}

type SignupForm struct {
    FirstName   string `schema:"firstName"`
    LastName    string `schema:"lastName"`
    Email       string `schema:"email"`
    Passwd      string `schema:"password"`
    SecretKey   string `schema:"secretkey"`
}

type LoginForm struct {
    Email       string `schema:"email"`
    Passwd      string `schema:"password"`
}

type Vaults struct {
    NewView     *views.View
    ShowView    *views.View
    vaultSrv    models.VaultService
    router      *mux.Router
}

type VaultForm struct {
    Email       string `schema:"email"`
    Username    string `schema:"username"`
    Application string `schema:"application"`
    Passwd      string `schema:"password"`
}

type Static struct {
    Home    *views.View
    Contact *views.View
}
