package models

import (
    "fmt"
    "log"
    "strings"
    "time"

    "github.com/jinzhu/gorm"

    "golang.org/x/crypto/ssh/terminal"
)

func UserInput(prompt string) (string, error) {
	fmt.Printf("%s: ", prompt)
	byte_input, err := terminal.ReadPassword(0)
	if err != nil {
        return "", err
	}
	user_input := string(byte_input)
    user_input = strings.TrimSpace(user_input)
    fmt.Println()

    return user_input, nil
}

/**
 * @brief:  Ask user for input and hide the input
 *
 * @param:  msg - Message to ask for user input
 *
 * @return: input on success, else error
 **/
func HiddenInput(msg string) string {
    attempts := 3
    for attempts > 0 {
        input_1, err := UserInput("Enter " + msg)
        if err != nil {
            log.Fatal(err)
        }

        input_2, err := UserInput("Re-enter " + msg)
        if err != nil {
            log.Fatal(err)
        }

        if input_1 == input_2 {
            return input_1
        }

        attempts--
        fmt.Printf("%s don't match\n", msg)
        time.Sleep(time.Second)
    }

    return ""
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
