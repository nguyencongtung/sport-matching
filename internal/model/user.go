package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// StringArray is a custom type for storing string arrays in the database as JSON.
type StringArray []string

// Value implements the driver.Valuer interface.
func (sa StringArray) Value() (driver.Value, error) {
	if sa == nil {
		return nil, nil
	}
	j, err := json.Marshal(sa)
	return string(j), err
}

// Scan implements the sql.Scanner interface.
func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, sa)
}

// User struct
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null;size:50;" validate:"required,min=3,max=50" json:"username"`
	Email    string `gorm:"uniqueIndex;not null;size:255;" validate:"required,email" json:"email"`
	Password string `gorm:"not null;" validate:"required,min=6,max=50" json:"password"`
	Names             string `json:"names"`
	Gender            string `json:"gender"`
	DateOfBirth       string `json:"date_of_birth"` // YYYY-MM-DD format
	Bio               string `gorm:"size:500;" json:"bio"`
	Interests         string `gorm:"size:255;" json:"interests"` // Comma-separated
	LookingFor        string `json:"looking_for"`
	ProfilePictureURLs StringArray `gorm:"type:jsonb" json:"profile_picture_urls"` // Store as JSONB in PostgreSQL
	Location          string `json:"location"`
	DistancePreference int    `json:"distance_preference"` // In kilometers
}
