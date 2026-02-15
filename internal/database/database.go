package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

type Case struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	VinNumber    string `json:"vin_number"`
	Approved     bool   `json:"approved"`
	Manufactured bool   `json:"manufactured"`
}

type JSONFileDB struct {
	file     *os.File
	filePath string
	db       dbStruct
	mu       sync.RWMutex
}

func NewJSONFileDB(filePath string) (*JSONFileDB, error) {
	db := &JSONFileDB{
		filePath: filePath,
		db:       make(dbStruct),
	}

	// Open or create the file
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	db.file = file

	// Read existing data if file is not empty
	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	if stat.Size() > 0 {
		if err := json.NewDecoder(file).Decode(&db.db); err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to decode JSON: %w", err)
		}
	}

	return db, nil
}

func (db *JSONFileDB) Save(c Case) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	id := uuid.NewString()
	db.db[id] = c

	if err := db.writeToFile(); err != nil {
		return "", err
	}

	return id, nil
}

func (db *JSONFileDB) Get(id string) (Case, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	c, ok := db.db[id]
	if !ok {
		return Case{}, fmt.Errorf("not found")
	}
	return c, nil
}

func (db *JSONFileDB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.file != nil {
		return db.file.Close()
	}
	return nil
}

func (db *JSONFileDB) writeToFile() error {
	// Truncate and seek to beginning
	if err := db.file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}
	if _, err := db.file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	// Write JSON
	encoder := json.NewEncoder(db.file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(db.db); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return db.file.Sync()
}

type dbStruct map[string]Case

type InMemoryDB struct {
	db dbStruct
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		db: make(dbStruct),
	}
}

func (db *InMemoryDB) Save(c Case) (string, error) {
	id := "abc-def"
	db.db[id] = c
	return id, nil
}

func (db *InMemoryDB) Get(id string) (Case, error) {
	c, ok := db.db[id]
	if !ok {
		return Case{}, fmt.Errorf("not found")
	}
	return c, nil
}

func (db *InMemoryDB) Close() error {
	return nil
}
