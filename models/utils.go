package models

import (
	"fmt"
	"strings"
    "time"

    "github.com/jinzhu/gorm"

	"golang.org/x/crypto/ssh/terminal"
)

/**
 * @brief:  Ask user for input and hide the input
 *
 * @param:  msg - Message to ask for user input
 *
 * @return: input on success, else error
 **/
func HiddenInput(msg string) string {
    for {
        fmt.Printf("Enter %s: ", msg)
        value, err := terminal.ReadPassword(0)
        if err != nil {
            panic(err)
        }
        input_1 := string(value)
        input_1 = strings.TrimSpace(input_1)
        fmt.Println("")

        fmt.Printf("Re-enter %s: ", msg)
        value, err = terminal.ReadPassword(0)
        if err != nil {
            panic(err)
        }
        input_2 := string(value)
        input_2 = strings.TrimSpace(input_2)
        fmt.Println("")

        if input_1 != input_2 {
            fmt.Printf("%s don't match\n", msg)
            time.Sleep(time.Second)
        } else {
            return input_1
        }
    }
}

/**
 * @brief:  Query provided database and get the first item returned
 *
 * @param:  db - Database to query
 * @param:  dst - Store returned data if found
 *
 * @return: nil on success, else error
 **/
func first(db *gorm.DB, dst interface{}) error {
    err := db.First(dst).Error
    if err == gorm.ErrRecordNotFound {
        return ErrNotFound
    }

    return err
}

/**
 * @brief:  Query provided database and find given items
 *
 * @param:  db - Database to query
 * @param:  dst - Store returned data if found
 *
 * @return: nil on success, else error
 **/
func find(db *gorm.DB, dst interface{}) error {
    err := db.Find(dst).Error
    if err == gorm.ErrRecordNotFound {
        return ErrNotFound
    }

    return err
}
