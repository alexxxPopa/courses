package sql

import (
	"github.com/jinzhu/gorm"
	"github.com/alexxxPopa/courses/conf"
	"github.com/pkg/errors"
	"github.com/alexxxPopa/courses/models"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Connection struct {
	db *gorm.DB
}

func (conn *Connection) Close() error {
	return conn.db.Close()
}

func (conn *Connection) Migrate() error {
	conn.db = conn.db.AutoMigrate(&models.User{}, &models.Adress{})
	return conn.db.Error
}

func Connect(config *conf.Config) (*Connection, error) {
	db, err := gorm.Open(config.DB.Driver, config.DB.Url)

	if err != nil {
		return nil, errors.Wrap(err, "Error while opening the database")
	}

	err = db.DB().Ping()

	if err != nil {
		return nil, errors.Wrap(err, "Error while connecting to the database")
	}

	conn := &Connection{
		db: db,
	}

	return conn, nil
}