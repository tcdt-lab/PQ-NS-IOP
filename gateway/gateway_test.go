package main

import (
	"database/sql"
	"testing"
)

//func TestBootLogic(t *testing.T) {
//	c, err := config.ReadYaml()
//	if err != nil {
//		t.Errorf("Error reading config.yaml file")
//	}
//	BootLogic(c)
//}
//func TestPhaseOneExecute(t *testing.T) {
//	PhaseOneExecute()
//}

func TestGetBalance(t *testing.T) {
	// Test the GetBalance function
	// You can use a mock database connection for testing
	// db := ...
	// balance, err := GetBalance(db)
	// if err != nil {
	// 	t.Errorf("Error getting balance: %v", err)
	// }
	// fmt.Printf("Balance: %d\n", balance)
	db, _ := getTestDbConnection()
	PhaseOneExecute(db)   // Replace with actual db connection
	BalanceCheckStart(db) // Replace with actual db connection
}

func getTestDbConnection() (*sql.DB, error) {
	c, err := getConfig()
	db, err := sql.Open("mysql", c.DB.Username+":"+c.DB.Password+"@/"+c.DB.Name)
	db.SetMaxOpenConns(1000)
	if err != nil {
		return nil, err
	}
	return db, nil
}
