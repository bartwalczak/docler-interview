package server

import (
	"fmt"

	"github.com/bartwalczak/docler-interview/api-go/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // just import a dialect
)

// DBConfig holds the config of the database connection
type DBConfig struct {
	Port string
	Host string
	User string
	Pass string
	Name string
}

// Open connects to the db
func (cfg DBConfig) Open() (*gorm.DB, error) {
	str := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Name)
	db, err := gorm.Open("postgres", str)
	if err != nil {
		return nil, err
	}
	// migrate database if needed
	db.AutoMigrate(models.Models()...)
	return db, nil
}
