package main

import (
    "fmt"
    "net/http"

    "github.com/loerac/vaultDepot/controllers"
    "github.com/loerac/vaultDepot/middleware"
    "github.com/loerac/vaultDepot/models"

    "github.com/gorilla/mux"
)

const (
    host    = "localhost"
    port    = 5432
    user    = "postgres"
    passwd  = "your-password"
    dbname  = "your-dbname"
)

func main() {
    /* Create DB connection */
    connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, passwd, dbname)

    services, err := models.NewServices(connInfo)
    if err != nil {
        panic(err)
    }
    defer services.Close()
    services.AutoMigrate()

    /* Mux Router */
    router := mux.NewRouter()

    /* Controllers */
    staticCtrl := controllers.NewStatic()
    usersCtrl := controllers.NewUsers(services.User)
    vaultCtrl := controllers.NewVaults(services.Vault, router)

    userMw := middleware.User {
        UserService: services.User,
    }
    reqUser := middleware.RequireUser{}
    newVault := reqUser.ApplyFn(vaultCtrl.New)
    createVault := reqUser.ApplyFn(vaultCtrl.Create)
    showVault := reqUser.ApplyFn(vaultCtrl.Show)

    /* Routes */
    router.Handle("/", staticCtrl.Home).Methods("GET")
    router.Handle("/contact", staticCtrl.Contact).Methods("GET")
    router.HandleFunc("/signup", usersCtrl.New).Methods("GET")
    router.HandleFunc("/signup", usersCtrl.Create).Methods("POST")
    router.Handle("/login", usersCtrl.LoginView).Methods("GET")
    router.HandleFunc("/login", usersCtrl.Login).Methods("POST")
    router.HandleFunc("/vault/new", newVault).Methods("GET")
    router.HandleFunc("/vault", createVault).Methods("POST")
    router.HandleFunc("/vault/account", showVault).Methods("GET")

    /* Dev testing: display cookies set on the current user */
    router.HandleFunc("/cookietest", usersCtrl.CookiesTest).Methods("GET")

    fmt.Println("Listening on http://localhost:3000")
    fmt.Println(http.ListenAndServe(":3000", userMw.Apply(router)))
}

