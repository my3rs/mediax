package database

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/scenery/mediax/helpers"
	"github.com/scenery/mediax/models"
	"github.com/scenery/mediax/version"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB() {
	var err error
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:        time.Second,
			LogLevel:             logger.Error,
			ParameterizedQueries: true,
			Colorful:             true,
		},
	)

	dbPath := "mediax.db"

	_, err = os.Stat(dbPath)

	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		TranslateError: true,
		Logger:         dbLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 当前版本低于 0.6.0
	if isOldVersion(version.Version, "0.6.0") {
		err = db.AutoMigrate(&models.Subject{})
		if err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}

		// 删除过时的索引
		if db.Migrator().HasIndex(&models.Subject{}, "idx_title") {
			db.Migrator().DropIndex(&models.Subject{}, "idx_title")
		}
		if db.Migrator().HasIndex(&models.Subject{}, "idx_id_type") {
			db.Migrator().DropIndex(&models.Subject{}, "idx_id_type")
		}
		if db.Migrator().HasIndex(&models.Subject{}, "idx_type_status_id") {
			db.Migrator().DropIndex(&models.Subject{}, "idx_type_status_id")
		}

		fmt.Println("Database migration successful.")
	}

}

func GetDB() *gorm.DB {
	return db
}

func isOldVersion(current, target string) bool {
	curParts := strings.Split(current, ".")
	tgtParts := strings.Split(target, ".")

	for len(curParts) < 3 {
		curParts = append(curParts, "0")
	}
	for len(tgtParts) < 3 {
		tgtParts = append(tgtParts, "0")
	}

	for i := 0; i < 3; i++ {
		cur, err1 := helpers.StringToInt(curParts[i])
		tgt, err2 := helpers.StringToInt(tgtParts[i])

		if err1 != nil || err2 != nil {
			return false
		}

		if cur < tgt {
			return true
		} else if cur > tgt {
			return false
		}
	}
	return false
}
