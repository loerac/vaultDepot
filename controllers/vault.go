package controllers

import (
    "net/http"

    "github.com/loerac/vaultDepot/context"
    "github.com/loerac/vaultDepot/models"
    "github.com/loerac/vaultDepot/views"

    "github.com/gorilla/mux"
)

const ShowVault = "show_vault"

/**
 * @brief:  Initialize the layout and template for vaults
 *
 * @return: Vaults controller
 **/
func NewVaults(vaultSrv models.VaultService, router *mux.Router) *Vaults {
    return &Vaults {
        NewView:    views.NewView("bootstrap", "vault/new"),
        ShowView:   views.NewView("bootstrap", "vault/show"),
        vaultSrv:   vaultSrv,
        router:     router,
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
func (vaultRW *Vaults) New(writer http.ResponseWriter, request *http.Request) {
    vaultRW.NewView.Render(writer, request, nil)
}

/**
 * @brief:  Process the signup form when a vault
 *          tries to create a new vault account
 *
 * @param:  writer - Render the html template
 * @param:  request - Input data from vault
 *
 * @action: POST /vault
 **/
func (vaultRW *Vaults) Create(writer http.ResponseWriter, request *http.Request) {
    viewData := views.Data{}
    form := VaultForm{}
    if err := parseForm(request, &form); err != nil {
        viewData.SetAlert(err)
        vaultRW.NewView.Render(writer, request, viewData)
        return
    }

    user := context.User(request.Context())
    vault := models.Vault {
        UserID:         user.ID,
        SecretKey:      user.SecretKey,
        Email:          form.Email,
        Username:       form.Username,
        Application:    form.Application,
        Password:       form.Passwd,
    }
    if err := vaultRW.vaultSrv.Create(&vault); err != nil {
        viewData.SetAlert(err)
        vaultRW.NewView.Render(writer, request, viewData)
        return
    }

    http.Redirect(writer, request, "/vault/account", http.StatusFound)
}

/**
 * @brief:  Display the item's in the vault
 *
 * @param:  writer - Render the html template
 * @param:  request - Input data from vault
 *
 * @action: GET /vault/account
 **/
func (vaultRW *Vaults) Show(writer http.ResponseWriter, request *http.Request) {
    key := ""
    user := context.User(request.Context())
    if user == nil {
        /* No user found, login */
        http.Redirect(writer, request, "/login", http.StatusFound)
        return
    }
    key = user.SecretKey
    vault, err := vaultRW.vaultSrv.ByID(uint(user.ID), key)
    if err != nil {
        switch err  {
        case models.ErrNotFound:
            http.Error(writer, "Vault not found", http.StatusNotFound)
        default:
            http.Error(writer, "OOPSIE WOOPSIE!! uwu something went vewy wrong owo", http.StatusInternalServerError)
        }
        return
    }

    type tempStr struct {
        Vault *[]models.Vault
    }
    dispVault := tempStr {
        Vault: vault,
    }

    viewData := views.Data{}
    viewData.Yield = dispVault
    vaultRW.ShowView.Render(writer, request, viewData)
}
