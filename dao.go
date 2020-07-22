package main

import (
	"log"

	"github.com/prologic/bitcask"
)

// Dao handles DB operations
type Dao struct {
	dbPath string
}

// DbOperationStatus represents the possibile outcome after accessing the DB
type DbOperationStatus int

// DbOperationStatus values
const (
	Ok DbOperationStatus = iota
	NotFound
	Error
)

// NewDao returns a dao that handles operations on the provided DB path
func NewDao(dbPath string) *Dao {
	return &Dao{dbPath}
}

// Save stores the value using the provided key. It returns Ok if saving is
// successful, Error otherwise
func (dao *Dao) Save(key, value string) DbOperationStatus {
	db, err := bitcask.Open(dao.dbPath)
	if err != nil {
		log.Printf("Unable to open db: %v", err)
		return Error
	}
	defer db.Close()

	err = db.Put([]byte(key), []byte(value))
	if err != nil {
		log.Printf("Unable to save value: %v", err)
		return Error
	}

	return Ok
}

// FindByKey returns the value, if found, and an Ok status. If not found or in
// case of errors, return empty string and the proper DBOperationStatus
func (dao *Dao) FindByKey(key string) (string, DbOperationStatus) {
	db, err := bitcask.Open(dao.dbPath)
	if err != nil {
		log.Printf("Unable to open db: %v", err)
		return "", Error
	}
	defer db.Close()

	url, err := db.Get([]byte(key))
	if err == bitcask.ErrKeyNotFound {
		log.Printf("Not Found for key '%v' - %v", key, err)
		return "", NotFound
	}

	if err != nil {
		log.Printf("Unable to retrieve value: %v", err)
		return "", Error
	}

	return string(url), Ok
}

// RemoveByKey deletes a value pointed by the provided key. Return Ok is
// deletion is successful, Error otherwise
func (dao *Dao) RemoveByKey(key string) DbOperationStatus {
	db, err := bitcask.Open(dao.dbPath)
	if err != nil {
		log.Printf("Unable to open db: %v", err)
		return Error
	}
	defer db.Close()

	byteKey := []byte(key)

	if !db.Has(byteKey) {
		return NotFound
	}

	err = db.Delete([]byte(key))

	if err != nil {
		log.Printf("Unable to delete value: %v", err)
		return Error
	}

	return Ok
}

// DoesExist check if the key is already present in the DB
func (dao *Dao) DoesExist(key string) (bool, DbOperationStatus) {
	db, err := bitcask.Open(dao.dbPath)
	if err != nil {
		log.Printf("Unable to open db: %v", err)
		return false, Error
	}
	defer db.Close()

	return db.Has([]byte(key)), Ok
}
