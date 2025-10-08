package connection

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Owner struct {
	gorm.Model

	Email string
	Phash string

	School School `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE;"`
}

type School struct {
	gorm.Model

	OwnerID uint

	Wireframes []Wireframe

	Users []User
}

type Wireframe struct {
	gorm.Model

	Name  string `gorm:"unique"`
	Data  datatypes.JSON
	Keys  datatypes.JSON
	Activ bool

	SchoolID uint
	School   School
}

type EmailToken struct {
	gorm.Model
	Email     string    `gorm:"primaryKey"`
	Code      string    `gorm:"not null;unique"`
	ExpiresAt time.Time `gorm:"not null"`
	Verified  bool      `gorm:"not null;default:false"`
}

type User struct {
	gorm.Model

	SchoolID uint
	School   School

	Role string

	Data  datatypes.JSON
	Email string `gorm:"not null;unique"`
	Phash string
}
