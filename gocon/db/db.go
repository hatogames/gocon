package connection

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Setup() {
	DB = connectDB(getDBurl())

	DB.AutoMigrate(&Owner{}, &School{}, &User{}, &Registration{}, &Student{}, &Wireframe{})
}

func getDBurl() string {
	dsn := "postgresql://neondb_owner:npg_9sVgw0DYTyoE@ep-divine-unit-ag4rv7l0-pooler.c-2.eu-central-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	return dsn
}

func connectDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	return db
}
