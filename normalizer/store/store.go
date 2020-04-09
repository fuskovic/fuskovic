package store

import (
	"os"
	"path"
	"unicode"

	"github.com/jinzhu/gorm"

	// sqlite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const dialect = "sqlite3"

var defaultNumbers = []string{
	"1234567890",
	"123 456 7891",
	"(123) 456 7892",
	"(123) 456-7893",
	"123-456-7894",
	"123-456-7890",
	"1234567892",
	"(123)456-7892",
}

type (
	// Store represents the persistence layer of this application.
	Store struct{ db *gorm.DB }

	entry struct {
		gorm.Model
		PhoneNumber string `gorm:"not null"`
	}

	// Entries contains multiple entries
	Entries []entry
)

// New creates a new store instance.
func New() (*Store, error) {
	db, err := gorm.Open(dialect, getBinPath())
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	if err := s.destructiveReset(); err != nil {
		return nil, err
	}
	return s, nil
}

// List lists all the phone numbers in the database
func (s *Store) List() (Entries, error) {
	var entries Entries
	if err := s.db.Find(&entries).Error; err != nil {
		return entries, err
	}
	return entries, nil
}

// Normalize normalizes all phone numbers in the store.
func (s *Store) Normalize() error {
	var entries, normalizedEntries Entries

	if err := s.db.Find(&entries).Error; err != nil {
		return err
	}

	for _, e := range entries {
		var normalizedNumber string
		for _, char := range e.PhoneNumber {
			if unicode.IsNumber(char) {
				normalizedNumber += string(char)
			}
		}
		e.PhoneNumber = normalizedNumber
		normalizedEntries = append(normalizedEntries, e)
	}

	counts := make(map[string]int)

	for _, e := range normalizedEntries {
		counts[e.PhoneNumber]++

		if counts[e.PhoneNumber] == 1 {
			if err := s.db.Where("id = ?", e.ID).Save(&e).Error; err != nil {
				return err
			}
		} else {
			if err := s.db.Where("id = ?", e.ID).Delete(&entry{}).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Store) destructiveReset() error {
	if err := s.db.DropTableIfExists(&entry{}).Error; err != nil {
		return err
	}

	if err := s.db.AutoMigrate(&entry{}).Error; err != nil {
		return err
	}

	for _, defaultNum := range defaultNumbers {
		if err := s.db.Create(&entry{PhoneNumber: defaultNum}).Error; err != nil {
			return err
		}
	}
	return nil
}

func getBinPath() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path.Join(wd, "numbers.db")
}
