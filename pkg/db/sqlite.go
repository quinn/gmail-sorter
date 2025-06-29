package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type OAuthAccount struct {
	ID        uint   `gorm:"primaryKey"`
	Provider  string `gorm:"index"`
	Email     string `gorm:"index"`
	TokenJSON string // serialized oauth2.Token
	Label     string // user-editable label (e.g. "Work", "Personal")
	CreatedAt int64  // unix timestamp
	UpdatedAt int64
}

// InitOAuthDB initializes the SQLite DB for OAuth accounts (singleton)
func InitOAuthDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&OAuthAccount{})
	return db, err
}

// CRUD helpers (optional, for clarity)
func CreateOAuthAccount(acct *OAuthAccount) error {
	db, err := InitOAuthDB()
	if err != nil {
		return err
	}
	return db.Create(acct).Error
}

func GetOAuthAccountByID(id uint) (*OAuthAccount, error) {
	db, err := InitOAuthDB()
	if err != nil {
		return nil, err
	}
	var acct OAuthAccount
	if err := db.First(&acct, id).Error; err != nil {
		return nil, err
	}
	return &acct, nil
}

func ListOAuthAccounts() ([]OAuthAccount, error) {
	db, err := InitOAuthDB()
	if err != nil {
		return nil, err
	}
	var accounts []OAuthAccount
	if err := db.Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func UpdateOAuthAccount(acct *OAuthAccount) error {
	db, err := InitOAuthDB()
	if err != nil {
		return err
	}
	return db.Save(acct).Error
}

func DeleteOAuthAccount(id uint) error {
	db, err := InitOAuthDB()
	if err != nil {
		return err
	}
	return db.Delete(&OAuthAccount{}, id).Error
}
