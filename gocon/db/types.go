package connection

import (
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

	Students      []Student
	Registrations []Registration
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

type Registration struct {
	gorm.Model
	Data datatypes.JSON

	SchoolID uint
	School   School

	UserID uint
	User   User
}

type Student struct {
	gorm.Model
	Data datatypes.JSON

	SchoolID uint
	School   School

	UserID uint
	User   User
}

type User struct {
	gorm.Model

	Email string
	Phash string

	Student      *Student      `gorm:"foreignKey:UserID"`
	Registration *Registration `gorm:"foreignKey:UserID"`
}
