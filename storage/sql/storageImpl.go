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

func (conn *Connection) FindUserByEmail(email string) (*models.User, error ){
	return conn.findUser("email = ?", email)
}

func (conn *Connection) CreateUser(user *models.User) error {
	tx := conn.db.Begin()
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (conn *Connection) findUser(query string, args ...interface{}) (*models.User, error) {
	user := &models.User{}
	values := append([]interface{}{query}, args...)
	if userExists := conn.db.First(user, values...); userExists.Error != nil {
		return nil, userExists.Error
	}
	return user, nil
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