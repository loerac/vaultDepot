package controllers

import (
    "net/http"

    "github.com/gorilla/schema"
)

/**
 * @brief:  Parse requested form
 *
 * @param:  request - HTTP requested data in a form
 * @param:  dst - Storing requested data
 *
 * @return: nil on success, else error
 **/
func parseForm(request *http.Request, dst interface{}) error {
    if err := request.ParseForm(); err != nil {
        return err
    }

    dec := schema.NewDecoder()
    if err := dec.Decode(dst, request.PostForm); err != nil {
        return err
    }

    return nil
}
