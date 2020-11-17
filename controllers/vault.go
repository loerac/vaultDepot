package controllers

import (
    "fmt"
    "net/http"

    "github.com/loerac/vaultDepot/models"
    "github.com/loerac/vaultDepot/views"
)

/**
 * @brief:  Initialize the layout and template for vaults
 *
 * @return: Vaults controller
 **/
func NewVaults(vaultSrv *models.VaultService) *Vaults {
    return &Vaults {
        VaultView: views.NewView("bootstrap", "vault/new"),
        vaultSrv: vaultSrv,
    }
}

/* ==============================*/
/*     METHODS FOR RESTful API   */
/* ==============================*/

/**
 * @brief:  Process the signup form when a vault
 *          tries to create a new vault account
 *
 * @param:  writer - Render the html template
 * @param:  request - Input data from vault
 *
 * @action: POST /vault
 **/
func (vaultRW *Vaults) Vault(writer http.ResponseWriter, request *http.Request) {
    form := VaultForm{}
    if err := parseForm(request, &form); err != nil {
        panic(err)
    }

    vault := models.Vault {
        Email:      form.Email,
        Username:   form.Username,
        Website:    form.Website,
        Password:   form.Passwd,
    }
    err := vaultRW.vaultSrv.Create(&vault)
    if err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintln(writer, "Created vault entry:", vault)
}
