package controllers

import (
    "github.com/loerac/vaultDepot/models"
    "github.com/loerac/vaultDepot/views"
)

type Users struct {
    NewView *views.View
    LoginView *views.View
    userSrv *models.UserService
}

type SignupForm struct {
    FirstName   string `schema:"firstName"`
    LastName    string `schema:"lastName"`
    Email       string `schema:"email"`
    Passwd      string `schema:"password"`
}

type LoginForm struct {
    Email       string `schema:"email"`
    Passwd      string `schema:"password"`
}

type Vaults struct {
    VaultView *views.View
    vaultSrv *models.VaultService
}

type VaultForm struct {
    Email       string `schema:"email"`
    Username    string `schema:"username"`
    Website     string `schema:"website"`
    Passwd      string `schema:"password"`
}

type Static struct {
    Home    *views.View
    Contact *views.View
}
