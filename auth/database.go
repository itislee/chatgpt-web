package main

import (
	"fmt"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Database interface {
	IsOpenIDExists(openID, accessToken string) (bool, error)
	UpdateAccessToken(openID, accessToken string) error
}

type MySQLDatabase struct {
	*sql.DB
}

func NewMySQLDatabase(dsn string) (*MySQLDatabase, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &MySQLDatabase{db}, nil
}

func (db *MySQLDatabase) IsOpenIDExists(openID, accessToken string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM auth_openid WHERE openid=? AND accesstoken=?", openID, accessToken).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (db *MySQLDatabase) UpdateAccessToken(openID, accessToken string) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	result, err := db.Exec("update auth_openid set accesstoken=?, updated_at=? where openid=?", accessToken, now, openID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Failed to update access token for openID: %s", openID)
	}

	return nil
}
