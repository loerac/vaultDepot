package models

import (
    "fmt"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
)

/**
 * @brief:  Connecting to vault database with GORM
 *
 * @param:  connInfo - Information of database
 *
 * @return: VaultService on success, else error
 **/
func NewVaultService(connInfo string) (*VaultService, error) {
    db, err := gorm.Open("postgres", connInfo)
    if err != nil {
        return nil, err
    }
    db.LogMode(true)

    return &VaultService{
        db: db,
    }, nil
}

/**
 * @brief:  Close the database connection
 *
 * @return: nil on success, else error
 **/
func (vaultSrv *VaultService) Close() error {
    return vaultSrv.db.Close()
}

/* ==============================*/
/*        METHODS FOR CRUD       */
/* ==============================*/

/**
 * @brief:  Create provided vault and backfill data
 *
 * @param:  vault - Vault information as struct
 *
 * @return: nil on success, else error
 **/
func (vaultSrv *VaultService) Create(vault *Vault) error {
    return vaultSrv.db.Create(vault).Error
}

/**
 * @brief:  Update provided vault with data
 *
 * @param:  vault - Vault information
 *
 * @return: nil on success, else error
 **/
func (vaultSrv *VaultService) Update(vault *Vault) error {
    return vaultSrv.db.Save(vault).Error
}

/**
 * @brief:  Delete provided vault with ID
 *
 * @param:  id - Vault to be deleted
 *
 * @return: nil on success
 *          ErrIDInvalid if ID is invalid
 *          Else error
 **/
func (vaultSrv *VaultService) Delete(id uint) error {
    if id == 0 {
        return ErrIDInvalid
    }

    vault := Vault{
        Model: gorm.Model {
            ID: id,
        },
    }
    return vaultSrv.db.Delete(&vault).Error
}

/* ==============================*/
/*   METHODS TO SEARCH FOR USER  */
/* ==============================*/

/**
 * @brief:  Look up a vault with provided ID.
 *
 * @param:  id  - ID of the vault
 *
 * @return: If vault is found, return nil
 *          If vault not found, return ErrNotFound
 *          Else, return error
 **/
func (vaultSrv *VaultService) ByID(id uint) (*Vault, error) {
    var vault Vault
    db := vaultSrv.db.Where("id = ?", id)
    err := first(db, &vault)
    if err != nil {
        return nil, err
    }

    return &vault, nil
}

/**
 * @brief:  Look up a vault with provided email addr
 *
 * @param:  email - Email address of vault
 *
 * @return: If vault is found, return nil
 *          If vault not found, return ErrNotFound
 *          Else, return error
 **/
func (vaultSrv *VaultService) ByEmail(email string) (*Vault, error) {
    var vault Vault
    db := vaultSrv.db.Where("email = ?", email)
    err := first(db, &vault)
    if err != nil {
        return nil, err
    }

    return &vault, nil
}

/**
 * @brief:  Authenticate a vault with provided
 *          email addr and password
 *
 * @param:  email - Vaults email address in DB
 * @param:  password - Vaults password for account
 *
 * @return: If email addr is invalid, return ErrNotFound
 *          If password is invalid, return ErrPasswordIncorrect
 *          If both are vaild, return vault
 *          Else, error
 **/
func (vaultSrv *VaultService) Authenticate(email, password string) (*Vault, error) {
    foundVault, err := vaultSrv.ByEmail(email)
    if err != nil {
        return nil, err
    }

    return foundVault, nil
}

/* ==============================*/
/*      METHODS FOR DATABASE     */
/*           MIGRATION           */
/* ==============================*/

/**
 * @brief:  Attempt to automatically migrate vault table
 *
 * @return: nil on success, else error
 **/
func (vaultSrv *VaultService) AutoMigrate() error {
    err := vaultSrv.db.AutoMigrate(&Vault{}).Error
    if err != nil {
        fmt.Println("models: Error on migrating Vault table: ", err)
    }

    return err
}

/**
 * @brief:  Drops the vault table and rebuilds it
 **/
func (vaultSrv *VaultService) DestructiveReset() error {
    err := vaultSrv.db.DropTableIfExists(&Vault{}).Error
    if err != nil {
        fmt.Println("models: Error on dropping Vault table: ", err)
        return err
    }

    return vaultSrv.AutoMigrate()
}
