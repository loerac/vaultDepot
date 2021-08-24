package main

import (
    "fmt"
    "strings"

    "github.com/atotto/clipboard"
    "github.com/loerac/vaultDepot/compat"
    "github.com/loerac/vaultDepot/models"
    "github.com/loerac/vaultDepot/manager"

    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
)

const (
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    password = "password"
    dbname   = "vaultdepot"
)

var checkError = compat.CheckError

var menu_options []string = []string{"Login", "Signup"}
var vault_menu_options []string = []string{"Get vault item", "Add vault item", "Export passwords", "Import passwords", "Exit"}
var vault_options []string = []string{"Copy password to clipboard", "Edit info", "Delete"}

var vaults []models.Vault

/**
 * @brief:  Display the options on the menu
 *
 * @arg:    options - Array with options on the menu
 *
 * @return: index of the option
 **/
func DisplayOptions(options []string) int {
    var input int
    for {
        for i, option := range options {
            fmt.Printf("%d.) %s\n", i + 1, option)
        }
        fmt.Printf("Choose (1 - %d): ", len(options))
        fmt.Scanln(&input)

        if input > 0 && input <= len(options) {
            break
        }
    }
    fmt.Println()

    return input
}

/**
 * @brief:  Display all the entries in the vault and select a specific vault
 *          entry
 *
 * @return: index of the vault selected
 **/
func getVaultItems() int {
    input := 0
    models.DisplayVault(vaults)
    for input != -1 {
        fmt.Print("Enter vault entry ('-1' to exit): ")
        fmt.Scanln(&input)
        if input < 0 || input > len(vaults) {
            if input != -1 {
                fmt.Println("Entry not found")
            }
            continue
        }
        break
    }
    fmt.Println()

    return input
}

/**
 * @brief:  Let user choose to copy password to clipboard, update the info in
 *          vault, or delete it.
 *
 * @arg:    db - Connection to the database
 * @arg:    user - User information to update a vault entry
 * @arg:    vault - Selected vault that was selected
 *
 * @return: index of the option
 **/
func selectVaultOptions(db *gorm.DB, user models.User, vault *models.Vault) {
    var char_input string

    fmt.Printf("Selected Vault: %v\n", *vault)
    input := DisplayOptions(vault_options)
    switch (input) {
    /* Copy password to clipboard */
    case 1:
        fmt.Println("Password copied to clipboard")
        clipboard.WriteAll(vault.Password)

    /* Edit entry */
    case 2:
        fmt.Printf("Update %s? (y or n): ", *vault)
        fmt.Scanln(&char_input)
        if strings.ToLower(char_input) == "y" {
            updated_vault, err := models.UpdateEntry(db, *vault, user)
            checkError(err)
            fmt.Printf("Updated: %s\n\n", updated_vault)
            vault = &updated_vault
        } else {
            fmt.Printf("\n%s wasn't updated\n\n", *vault)
        }

    /* Delete entry */
    case 3:
        fmt.Printf("Delete %s from vault? (y or n): ", *vault)
        fmt.Scanln(&char_input)
        if strings.ToLower(char_input) == "y" {
            err := models.DeleteID(db, vault.ID)
            checkError(err)
            fmt.Printf("%s was deleted\n\n", *vault)
        } else {
            fmt.Printf("\n%s wasn't deleted\n\n", *vault)
        }
    default:
        break
    }
}

func main() {
    /**
     * Log into the vault_database,
     * Enable logging,
     * Migrate the user, and vault tables if not created
     **/
    psqlinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    db, err := gorm.Open("postgres", psqlinfo)
    checkError(err)
    defer db.Close()
    db.LogMode(false)
    db.AutoMigrate(&models.User{}, &models.Vault{})

    var user models.User

    /* Let the user login or signup */
    menu := DisplayOptions(menu_options)
    if menu == 1 {
        user, err = models.Login(db)
    } else {
        user, err = models.Signup(db)
    }
    checkError(err)

    /**
     * Grab all the items from user's vault.
     * If nothing is in the vault, let them add it in
     **/
    fmt.Println("\nWelcome", user.Username)
    vaults, err = models.FindAll(db, user.ID)
    if err == models.ErrNotFound || len(vaults) == 0 {
        fmt.Println("Looks like you have nothing in your vault, let's update that")
        vault, err := models.CreateEntry(db, user)
        checkError(err)
        vaults = append(vaults, vault)
    } else if err != nil {
        panic(err)
    }

    /**
     * Main loop
     **/
    var input int
    for input != -1 {
        /**
         * Let user choose to get an item from the vaault,
         * add item to the vault, or exit
         **/
        input = DisplayOptions(vault_menu_options)
        switch (input) {
        /* Get vault item */
        case 1:
            input = getVaultItems()
            if -1 == input {
                break
            }

            entry := input - 1
            vaults[entry], err = models.ByID(db, vaults[entry].ID, user)
            checkError(err)

            selectVaultOptions(db, user, &vaults[entry])
            vaults, err = models.FindAll(db, user.ID)
            checkError(err)

        /* Add vault item */
        case 2:
            _, err = models.CreateEntry(db, user)
            checkError(err)

            vaults, err = models.FindAll(db, user.ID)
            checkError(err)

        /* Export vault */
        case 3:
            err = manager.ExportManager(vaults, user)
            checkError(err)

        /* Import vault */
        case 4:
            err = manager.ImportManager(db, user)
            checkError(err)

            vaults, err = models.FindAll(db, user.ID)
            checkError(err)

        default:
            fmt.Println("Bye bye")
            input = -1
        }
    }
}
