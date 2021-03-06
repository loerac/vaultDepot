package main

import (
	"fmt"

	"github.com/atotto/clipboard"
    "github.com/loerac/vaultDepot/models"
    "github.com/loerac/vaultDepot/compat"

    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
)

const (
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    password = "password"
    dbname   = "vault_depot"
)

var checkError = compat.CheckError

var menu_options []string = []string{"Login", "Signup"}
var vault_menu_options []string = []string{"Get vault item", "Add vault item", "Exit"}
var vault_options []string = []string{"Copy password to clipboard", "Edit info", "Delete"}

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
    db.LogMode(true)
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
    vaults, err := models.FindAll(db, user.ID)
    if err == models.ErrNotFound || len(vaults) == 0 {
        fmt.Println("Looks like you have nothing in your vault, let's update that")
        vault, err := models.VaultEntry(db, user)
        checkError(err)
        vaults = append(vaults, vault)
    } else if err != nil {
        panic(err)
    }

    /**
     * Main loop
     *
     * TODO: Make prettier
     **/
    var input int
    var char_input string
    for input != -1 {
        /**
         * Let user choose to get an item from the vaault,
         * add item to the vault, or exit
         **/
        input = DisplayOptions(vault_menu_options)
        switch (input) {
        case 1:
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
            if input == -1 {
                input = 0
                break
            }

            entry := input - 1
            vault, err := models.ByID(db, vaults[entry].ID, user)
            checkError(err)

            /**
             * Let user choose to copy password to clipboard,
             * update the info in vault, or delete it.
             **/
            fmt.Println("Found:", vault)
            input = DisplayOptions(vault_options)
            switch (input) {
            case 1:
                fmt.Println("Password copied to clipboard")
	            clipboard.WriteAll(vault.Password)
            case 2:
                fmt.Printf("Update %s from vault? (Y or n): ", vault)
                fmt.Scanln(&char_input)
                if char_input == "Y" {
                    updated_vault, err := models.UpdateEntry(db, vault)
                    checkError(err)
                    fmt.Printf("Updated: %s\n\n", updated_vault)
                    vaults[entry] = updated_vault
                } else {
                    fmt.Printf("\n%s wasn't updated\n\n", vault)
                }
            case 3:
                fmt.Printf("Delete %s from vault? (Y or n): ", vault)
                fmt.Scanln(&char_input)
                if char_input == "Y" {
                    err = models.DeleteID(db, vault.ID)
                    checkError(err)
                    fmt.Printf("%s was deleted\n\n", vault)
                } else {
                    fmt.Printf("\n%s wasn't deleted\n\n", vault)
                }
            default:
                break
            }
        case 2:
            _, err = models.VaultEntry(db, user)
            checkError(err)

            vaults, err = models.FindAll(db, user.ID)
            checkError(err)
        default:
            fmt.Println("Bye bye")
            input = -1
        }
    }
}
