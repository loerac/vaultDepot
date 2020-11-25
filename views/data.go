package views

import (
    "log"
)

const (
    AlertError  = "danger"
    ALertWarn   = "warning"
    AlertInfo   = "info"
    AlertSucc   = "success"

    AlertGeneric = "OOPSIE WOOPSIE!! uwu I made a fucky wucky... owo A wittle fucko boingo!"
)

type PublicError interface {
    error
    Public() string
}

type Alert struct {
    Level   string
    Message string
}

type Data struct {
    Alert   *Alert
    Yield   interface{}
}

/***
 * @brief:  Set the level and message of the alert
 *
 * @param:  level - Alert level (success, info, etc)
 * @param:  msg - Message of the alert
 ***/
func (data *Data) SetupAlert(level, msg string) {
    data.Alert = &Alert {
        Level:  level,
        Message: msg,
    }
}

/***
 * @brief:  Set up the alert with either the public or generic error
 *          A generic error is only given if the error has sensitive info
 *
 * @param:  err - Details of the error
 ***/
func (data *Data) SetAlert(err error) {
    var msg string

    if pErr, ok := err.(PublicError); ok {
        msg = pErr.Public()
    } else {
        log.Println(err)
        msg = AlertGeneric
    }

    data.Alert = &Alert {
        Level:  AlertError,
        Message: msg,
    }
}
