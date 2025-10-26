package connection

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Setup() {
	DB = connectDB(getDBurl())

	DB.AutoMigrate(&Owner{}, &School{}, &Registration{}, &Wireframe{}, &EmailToken{})
}

func getDBurl() string {
	dsn := "postgresql://neondb_owner:npg_9sVgw0DYTyoE@ep-divine-unit-ag4rv7l0-pooler.c-2.eu-central-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	return dsn
}

func connectDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger konfigurieren: SQL-Logs sichtbar, Warnungen wie 'record not found' unterdrücken
		Logger: logger.New(
			log.Default(),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Error, // nur echte Fehler loggen
				IgnoreRecordNotFoundError: true,         // <-- Wichtig!
				Colorful:                  true,
			},
		),
	})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	return db
}
