package main

import "fmt"

type Case struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Email     string `json:"email"`
	VinNumber string `json:"vin_number"`
}

// TODO: JSON file implementation for persistent storage

type dbStruct map[string]Case

type inMemoryDB struct {
	db dbStruct
}

func NewInMemoryDB() *inMemoryDB {
	return &inMemoryDB{
		db: make(dbStruct),
	}
}

func (db *inMemoryDB) Save(c Case) (string, error) {
	id := "abc-def"
	db.db[id] = c
	return id, nil
}

func (db *inMemoryDB) Get(id string) (Case, error) {
	c, ok := db.db[id]
	if !ok {
		return Case{}, fmt.Errorf("not found")
	}
	return c, nil
}
