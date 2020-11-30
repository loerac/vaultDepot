package models

import (
    "github.com/jinzhu/gorm"
)

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
