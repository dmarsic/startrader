package store

import (
	"github.com/abhinav-TB/dantdb"
)

const StoreDir = "./data/"

func ReadAllUsers() ([]string, error) {
	db, err := dantdb.New(StoreDir)
	if err != nil {
		return nil, err
	}

	records, err := db.ReadAll("users")
	if err != nil {
		return nil, err
	}

	return records, nil
}
