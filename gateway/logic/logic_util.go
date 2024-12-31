package logic

import (
	"database/sql"
	"gateway/config"
)

func getConfig() (config.Config, error) {
	cfg, err := config.ReadYaml()
	return *cfg, err
}

func getDbConnection() (*sql.DB, error) {
	c, err := getConfig()
	db, err := sql.Open("mysql", c.DB.Username+":"+c.DB.Password+"@/"+c.DB.Name)
	if err != nil {
		return nil, err
	}
	return db, nil
}
