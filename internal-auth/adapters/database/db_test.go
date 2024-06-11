package database

import (
	_ "embed"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

//
//func TestDB(t *testing.T) {
//
//	dbConn, ncErr := NewConnection(&config.Config{
//		DatabasePath: ":memory:",
//		DatabaseDSN:  "sqlite3://:memory:"})
//	assert.NoError(t, ncErr)
//
//	defer func(dbConn *DB) {
//		cErr := dbConn.Close()
//		assert.NoError(t, cErr)
//	}(dbConn)
//
//	pErr := dbConn.Ping()
//	assert.NoError(t, pErr)
//
//	err := dbConn.MakeMigrations()
//	assert.NoError(t, err)
//
//	cErr := dbConn.Close()
//	assert.NoError(t, cErr)
//}
