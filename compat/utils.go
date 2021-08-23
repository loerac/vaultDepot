package compat

import (
    "log"
    "strconv"
)

/**
 * @brief:  Converts a string to an uint64.
 *          Panic if errors
 *
 * @arg:    str - String that is to be converted to uint64
 *
 * @return: uint64 val of string
 **/
func StrToInt(str string) int64 {
    num, err := strconv.ParseInt(str, 10, 64)
    CheckError(err)

    return num
}

/**
 * @brief:  Check if any errors occured, panic if so
 *
 * @arg:    e - Error
 **/
func CheckError(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
