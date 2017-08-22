package sql

import (
	"github.com/jinzhu/gorm"
	"github.com/alexxxPopa/courses/conf"
	"github.com/pkg/errors"
	"github.com/alexxxPopa/courses/models"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/alexxxPopa/courses/api"
)

type Connection struct {
	db *gorm.DB
}

func (conn *Connection) Close() error {
	return conn.db.Close()
}

func (conn *Connection) Migrate() error {
	//conn.db.DropTable(&models.User{}, &models.Plan{}, &models.Subscription{})
	conn.db = conn.db.AutoMigrate(&models.User{}, &models.Plan{}, &models.Subscription{})
	return conn.db.Error
}

func (conn *Connection) FindPlans() ([]api.PlanInfo, error) {
	var plans [] api.PlanInfo

	rows, _ := conn.db.Model(&models.Plan{}).
		Select("tile, amount").Rows()

	for rows.Next() {
		plan := api.PlanInfo{}
		err := rows.Scan(plan)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}

	return plans, nil
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
