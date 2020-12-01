package models

import (
    "github.com/loerac/vaultDepot/compat"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
)

/**
 * @brief:  Initialize vault valitadtor and GORM with
 *          the database pointer
 *
 * @param:  db - GORM database pointer
 *
 * @return: New VaultService
 **/
func NewVaultService(db *gorm.DB) VaultService {
    return &vaultService {
        VaultDB:    &vaultValidator {
            VaultDB:    &vaultGorm {
                db: db,
            },
        },
    }
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
func (vaultgorm *vaultGorm) Create(vault *Vault) error {
    return vaultgorm.db.Create(vault).Error
}

/* ==============================*/
/*   METHODS TO SEARCH FOR USER  */
/* ==============================*/

/**
 * @brief:  Look up a user with provided ID.
 *
 * @param:  id  - ID of the user
 * @param:  key - Users key for the vault
 *
 * @return: If user is found, return nil
 *          If user not found, return ErrNotFound
 *          Else, return error
 **/
func (vaultgorm *vaultGorm) ByID(id uint, key string) (*[]Vault, error) {
    vault := []Vault{}
    db := vaultgorm.db.Where("user_id = ?", id)
    err := find(db, &vault)
    if err != nil {
        return nil, err
    }

    if key != "" {
        aes := compat.NewAES(key)
        for i := range vault {
            password, err := aes.Decrypt(vault[i].PasswordCipher)
            if err != nil {
                return nil, err
            }

            vault[i].Password = password
        }
    }

    return &vault, nil
}

func (vaultgorm *vaultGorm) ByUserID(userID uint) ([]Vault, error) {
    vault := []Vault{}
    db := vaultgorm.db.Where("user_id = ?", userID)
    if err := db.Find(&vault).Error; err != nil {
        return nil, err
    }

    return vault, nil
}
