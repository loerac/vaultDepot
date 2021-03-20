package manager

import (
    "encoding/csv"
    "fmt"
    "os"

    "github.com/loerac/vaultDepot/compat"
    "github.com/loerac/vaultDepot/models"
    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
)

/**
 * CSV row header
 **/
var header = []string{"email", "username", "application", "password"}

/**
 * @brief:  Ask the user for a filename path
 *
 * @param:  port_type - Importing (true) or exporting (false) a CSV file
 *
 * @return: Filename path
 **/
func Filename(port_type bool) string {
    _type := "Exporting"
    if port_type {
        _type = "Importing"
    }

    filename := ""
    for filename == "" {
        fmt.Print("Enter filename: ")
        fmt.Scanln(&filename)
    }

    ext := filename[len(filename) - 4:]
    if ".csv" != ext {
        filename += ".csv"
    }

    fmt.Printf("%s %s...\n", _type, filename)
    return filename
}

/**
 * @brief:  Import a CSV with the header values (email, username, application, password)
 *          Password is textbase, and will be encrypted with users secret key
 *
 * @param:  db - pointer to dabase
 * @param:  user - contains user ID
 *
 * @return: nil on success, else error
 **/
//func ImportManager(db *gorm.DB, user models.User) ([]models.Vault, error) {
func ImportManager(db *gorm.DB, user models.User) error {
    filename := Filename(true)
    file, err := os.Open(filename)
    if nil != err {
        //return []models.Vault{}, err
        return err
    }
    defer file.Close()

    read := csv.NewReader(file)
    rows, err := read.ReadAll()
    if err != nil {
        fmt.Printf("Failed to import %s: %s.\n", filename, err)
        return err
    }

    vaults := []models.Vault{}
    for _, row := range rows[1:] {
        vault := models.Vault{
            UserID: user.ID,
            Email: row[0],
            Username: row[1],
            Application: row[2],
            Password: row[3],
        }

        vault, err := models.VaultEntry(db, vault, user)
        if nil != err {
            //return []models.Vault{}, err
            return err
        }
        vaults = append(vaults, vault)
    }

    fmt.Printf("Imported %s to vault\n", filename)

    //return vaults, nil
    return nil
}

/**
 * @brief:  Export a CSV with. Password will be in textbase.
 *
 * @param:  vaults - entries that will be exported
 * @param:  user - decrypt password
 *
 * @return: nil on success, else error
 **/
func ExportManager(vaults []models.Vault, user models.User) error {
    filename := Filename(false)
    file, err := os.Create(filename)
    if nil != err {
        return err
    }
    defer file.Close()

    write := csv.NewWriter(file)
    defer write.Flush()

    err = write.Write(header)
    if nil != err {
        return err
    }

    for _, vault := range vaults {
        aes := compat.NewAES(user.SecretKey)
        password, err := aes.Decrypt(vault.PasswordCipher)
        if err != nil {
            fmt.Printf("Failed to decrypt %s, skipping...\n", vault)
            continue
        }
        vault.Password = password

        var data = []string{vault.Email,
                            vault.Username,
                            vault.Application,
                            vault.Password,
                           }
        err = write.Write(data)
        if nil != err {
            return err
        }
    }

    fmt.Printf("Exported vault to %s\n", filename)

    return nil
}
