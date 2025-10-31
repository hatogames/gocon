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

	Registrations []Registration
}

type Wireframe struct {
	gorm.Model

	Name  string `gorm:"not null;unique"`
	Data  datatypes.JSON
	Keys  datatypes.JSON
	Activ bool

	SchoolID uint
	School   School
}

type EmailToken struct {
	gorm.Model
	Email     string    `gorm:"not null;unique"`
	Code      string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

type Registration struct {
	gorm.Model

	SchoolID uint
	School   School

	WireframeName string

	Data  datatypes.JSON
	Email string `gorm:"not null;unique"`
}
