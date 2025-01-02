package data_access

import (
	"database/sql"
	"gateway/config"
	"testing"
)

func TestLogin(t *testing.T) {

	//c, err := config.ReadYaml()
	//if err != nil {
	//	t.Errorf("Error reading config.yaml file: %v", err)
	//}
	//Login(getDBConnection(*c))

}

func getDBConnection(c config.Config) *sql.DB {
	db, err := sql.Open("mysql", c.DB.Username+":"+c.DB.Password+"@/"+c.DB.Name)
	if err != nil {
		return nil
	}
	return db
}
