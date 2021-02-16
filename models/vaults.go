package models

import (
	"bufio"
	"fmt"
	"os"
	"strings"

    "github.com/loerac/vaultDepot/compat"
    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
)

/**
 * @brief:  Ask for info to add to the vault
 *
 * @param:  id - User's ID to map item to user
 **/
func EntryInfo(id uint) Vault {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
    email = strings.TrimSpace(email)

	fmt.Print("Enter username (optional): ")
	username, _ := reader.ReadString('\n')
    username = strings.TrimSpace(username)

    fmt.Print("Enter application: ")
    app, _ := reader.ReadString('\n')
    app = strings.TrimSpace(app)

    password := HiddenInput("password")

    return Vault {
        UserID: id,
        Email: email,
        Username: username,
        Application: app,
        Password: password,
    }
}

/**
 * @brief:  Display all the items on the vault; only the application and email
 *
 * @param:  vaults - Array of item in the vault
 **/
func DisplayVault(vaults []Vault) {
    for i, vault := range vaults {
        fmt.Printf("%d.)\n\tApplication: %s\n\tEmail: %s\n", i + 1, vault.Application, vault.Email)
    }
}

/**
 * @brief:  Add an item from the vault
 *
 * @param:  db - pointer to dabase
 * @param:  user - contains user ID
 *
 * @return: New vault on success, else error
 **/
func VaultEntry(db *gorm.DB, user User) (Vault, error) {
    vault := EntryInfo(user.ID)

    if err := CreateEntry(db, &vault, user); err != nil {
        return Vault{}, err
    }

    fmt.Println("New item added to your vault:", vault)
    fmt.Println()

	return vault, nil
}

/**
 * @brief:  Update an item from the vault
 *
 * @param:  db - pointer to dabase
 * @param:  vault - vault to update
 *
 * @return: Updated vault on success, else error
 **/
func UpdateEntry(db *gorm.DB, vault Vault) (Vault, error) {
    updated_vault := EntryInfo(vault.UserID)
    updated_vault.ID = vault.ID

    if err := db.Save(&updated_vault).Error; err != nil {
        return Vault{}, err
    }

	return updated_vault, nil
}

/**
 * @brief:  Find first vaults with provided ID.
 *
 * @param:  id  - ID of the vault
 *
 * @return: If vault is found, return vault
 *          If vault not found, return ErrNotFound
 *          Else, return error
 **/
func ByID(vaultdb *gorm.DB, id uint, user User) (Vault, error) {
    var vault Vault
    db := vaultdb.Where("id = ?", id)
    err := first(db, &vault)
    if err != nil {
        return Vault{}, err
    }

    aes := compat.NewAES(user.SecretKey)
    password, err := aes.Decrypt(vault.PasswordCipher)
    if err != nil {
        return Vault{}, err
    }

    vault.Password = password

    return vault, nil
}

/**
 * @brief:  Find all vaults with provided ID.
 *
 * @param:  vaultdb - pointer to database
 * @param:  id  - ID of the user
 *
 * @return: If user is found, return vault
 *          If user not found, return ErrNotFound
 *          Else, return error
 **/
func FindAll(vaultdb *gorm.DB, id uint) ([]Vault, error) {
    var vault []Vault
    db := vaultdb.Where("user_id = ?", id)
    err := find(db, &vault)
    if err != nil {
        return nil, err
    }

    return vault, nil
}

/**
 * @brief:  Removes item from vault
 *
 * @param:  db - pointer to dabase
 * @param:  id - ID of item in vault
 *
 * @return: nil on success, else error
 **/
func DeleteID(vaultdb *gorm.DB, id uint) error {
    vault := Vault {
        Model: gorm.Model {
            ID: id,
        },
    }

    return vaultdb.Delete(&vault).Error
}
