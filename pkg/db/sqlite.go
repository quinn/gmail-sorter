package db

import "gorm.io/gorm"

type OAuthAccount struct {
	ID        uint   `gorm:"primaryKey"`
	Provider  string `gorm:"index"`
	Email     string `gorm:"index"`
	TokenJSON string // serialized oauth2.Token
	Label     string // user-editable label (e.g. "Work", "Personal")
	CreatedAt int64  // unix timestamp
	UpdatedAt int64
}

// CreateOAuthAccount inserts a new OAuthAccount
func (d *DB) CreateOAuthAccount(acct *OAuthAccount) error {
	return d.gorm.Create(acct).Error
}

// GetOAuthAccountByID retrieves an OAuthAccount by ID
func (d *DB) GetOAuthAccountByID(id string) (*OAuthAccount, error) {
	var acct OAuthAccount
	if err := d.gorm.First(&acct, id).Error; err != nil {
		return nil, err
	}
	return &acct, nil
}

// UpsertOAuthAccount inserts or updates an OAuthAccount by provider and email
func (d *DB) UpsertOAuthAccount(acct *OAuthAccount) error {
	var existing OAuthAccount
	err := d.gorm.Where("provider = ? AND email = ?", acct.Provider, acct.Email).First(&existing).Error
	if err == nil {
		acct.ID = existing.ID
		return d.gorm.Save(acct).Error
	} else if err.Error() == "record not found" || err == gorm.ErrRecordNotFound {
		return d.gorm.Create(acct).Error
	} else {
		return err
	}
}

// ListOAuthAccounts returns all OAuthAccounts
func (d *DB) ListOAuthAccounts() ([]OAuthAccount, error) {
	var accounts []OAuthAccount
	if err := d.gorm.Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

// UpdateOAuthAccount updates an existing OAuthAccount
func (d *DB) UpdateOAuthAccount(acct *OAuthAccount) error {
	return d.gorm.Save(acct).Error
}

// DeleteOAuthAccount deletes an OAuthAccount by ID
func (d *DB) DeleteOAuthAccount(id string) error {
	return d.gorm.Delete(&OAuthAccount{}, id).Error
}
