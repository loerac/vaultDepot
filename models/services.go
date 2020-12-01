package models

import (
    "log"

    "github.com/jinzhu/gorm"
)

/**
 * @brief:  Set up user and vault database with GORM
 *
 * @param:  connInfo - Information of database
 *
 * @return: Services on success, else error
 **/
func NewServices(connInfo string) (*Services, error) {
    db, err := gorm.Open("postgres", connInfo)
    if err != nil {
        return nil, err
    }
    db.LogMode(true)

    return &Services {
        db:     db,
        User:   NewUserService(db),
        Vault:  NewVaultService(db),
    }, nil
}

/**
 * @brief:  Close the database connection
 *
 * @return: error if closing database failed
 **/
func (srv *Services) Close() error {
    return srv.db.Close()
}

/* ==============================*/
/*      METHODS FOR DATABASE     */
/*           MIGRATION           */
/* ==============================*/

/**
 * @brief:  Attempt to automatically migrate all tables
 *
 * @return: error if migrating failed
 **/
func (srv *Services) AutoMigrate() error {
    err := srv.db.AutoMigrate(&User{}, &Vault{}).Error
    if err != nil {
        log.Println("models: Error on migrating user table:", err)
    }

    return err
}

/**
 * @brief:  Drop all tables and rebuild them
 *
 * @return: error if dropping table or migrating failed
 **/
func (srv *Services) DesctruciveReset() error {
    err := srv.db.DropTableIfExists(&User{}, &Vault{}).Error
    if err != nil {
        log.Println("models: Error on dropping user table:", err)
        return err
    }

    return srv.AutoMigrate()
}
