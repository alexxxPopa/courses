package sql

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/alexxxPopa/courses/conf"
	"github.com/stretchr/testify/suite"
	"github.com/alexxxPopa/courses/storage/sql/test"
	"github.com/alexxxPopa/courses/models"
)

var conn *Connection

func TestStorageTestSuite(t *testing.T) {
	config, err := conf.LoadTestConfig("../../config_test.json")
	require.NoError(t, err)

	conn, err = Connect(config)
	require.NoError(t, err)

	conn.db.Begin()
	defer conn.db.Callback()
	s := &test.StorageTestSuite{
		Conn:          conn,
		BeforeTest: beforeTest,
	}
	suite.Run(t, s)
}

func beforeTest() {
	conn.db.DropTableIfExists(&models.User{})
	conn.db.DropTableIfExists(&models.Plan{})
	conn.db.DropTableIfExists(&models.Subscription{})
	conn.db.DropTableIfExists(&models.Course{})
	conn.db.DropTableIfExists(&models.Article{})
	conn.db.DropTableIfExists(&models.Category{})
	conn.db.CreateTable(&models.User{})
	conn.db.CreateTable(&models.Plan{})
	conn.db.CreateTable(&models.Subscription{})
	conn.db.CreateTable(&models.Course{})
	conn.db.CreateTable(&models.Article{})
	conn.db.CreateTable(&models.Category{})
}