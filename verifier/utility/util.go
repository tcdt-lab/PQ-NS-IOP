package utility

import (
	"database/sql"
	"verifier/config"
)

func GetDBConnection(c config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", c.DB.Username+":"+c.DB.Password+"@/"+c.DB.Name)
	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(1000)
	if err != nil {
		return nil, err
	}
	return db, nil
}
