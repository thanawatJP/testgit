package database

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func Connect() {
	// โหลดไฟล์ .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// อ่านค่าจาก .env
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=" + os.Getenv("DB_SSLMODE")

	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // output logger
		logger.Config{
			SlowThreshold: time.Second, // ช่วงเวลาที่ถือว่าช้า
			LogLevel:      logger.Info, // ระดับการ log เช่น Info, Warn, Error
			Colorful:      true,        // ใช้สีใน Log
		},
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
}
