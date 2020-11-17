package controllers

import (
    "github.com/loerac/vaultDepot/views"
)

func NewStatic() *Static {
    return &Static {
        Home: views.NewView (
            "bootstrap", "static/home"),
        Contact: views.NewView (
            "bootstrap", "static/contact"),
    }
}
