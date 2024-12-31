package main

import (
	"database/sql"
	"fmt"
	"gateway/config"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"os"
)

func main() {
	c, err := config.ReadYaml()
	if err != nil {
		zap.L().Error("Error reading config.yaml file", zap.Error(err))
		os.Exit(1)
	}
	logger := zap.NewExample()
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()
	db, err := getDBConnection(*c)

	if err != nil {
		zap.L().Error("Error opening database", zap.Error(err))
		os.Exit(1)
	}
	defer db.Close()
	zap.L().Info("replaced zap's global loggers")
	fmt.Println("Welcome to the Gateway")

	//logic.Login(db)

}

func getDBConnection(c config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", c.DB.Username+":"+c.DB.Password+"@/"+c.DB.Name)
	if err != nil {
		return nil, err
	}
	return db, nil
}
