package main

import (
    "fmt"
    "net/http"

    "github.com/loerac/vaultDepot/controllers"
    "github.com/loerac/vaultDepot/models"

    "github.com/gorilla/mux"
)

const (
    host    = "localhost"
    port    = 5432
    user    = "postgres"
    passwd  = "your-password"
    dbname  = "database-name"
)

func main() {
    /* Create DB connection */
    connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, passwd, dbname)

    // TODO: Ugly, why two???
    userSrv, err := models.NewUserService(connInfo)
    if err != nil  {
        panic(err)
    }
    defer userSrv.Close()
    userSrv.AutoMigrate()

    vaultSrv, err := models.NewVaultService(connInfo)
    if err != nil  {
        panic(err)
    }
    defer vaultSrv.Close()
    vaultSrv.AutoMigrate()

    /* Controllers */
    staticCtrl := controllers.NewStatic()
    usersCtrl := controllers.NewUsers(userSrv)
    vaultCtrl := controllers.NewVaults(vaultSrv)

    /* Routes */
    router := mux.NewRouter()
    router.Handle("/", staticCtrl.Home).Methods("GET")
    router.Handle("/contact", staticCtrl.Contact).Methods("GET")
    router.HandleFunc("/signup", usersCtrl.NewSignup).Methods("GET")
    router.HandleFunc("/signup", usersCtrl.CreateSignup).Methods("POST")
    router.Handle("/login", usersCtrl.LoginView).Methods("GET")
    router.HandleFunc("/login", usersCtrl.Login).Methods("POST")
    router.Handle("/vault", vaultCtrl.VaultView).Methods("GET")
    router.HandleFunc("/vault", vaultCtrl.Vault).Methods("POST")

    fmt.Println("Listening on http://localhost:3000")
    fmt.Println(http.ListenAndServe(":3000", router))
}

