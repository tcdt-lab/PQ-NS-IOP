package data

import (
	"database/sql"
	"go.uber.org/zap"
)

type User struct {
	Id        int
	Password  string
	PublicKey string
	SecretKey string
	Salt      string
}

func GetUser(db *sql.DB, id int) (User, error) {
	var user User
	rows := db.QueryRow("SELECT * FROM gateway_user WHERE Id = ?", id)
	if err := rows.Scan(&user.Id, &user.Password, &user.PublicKey, &user.SecretKey, &user.Salt); err != nil {
		zap.L().Error("Error getting user", zap.Error(err))
		return User{}, err
	}
	return user, nil
}

func AddUser(db *sql.DB, user User) (int64, error) {
	result, err := db.Exec("INSERT INTO gateway_user (Password, PublicKey, secret_key, Salt) VALUES (?, ?, ?, ?)", user.Password, user.PublicKey, user.SecretKey, user.Salt)
	if err != nil {
		zap.L().Error("Error adding user", zap.Error(err))
		return 0, err
	}
	return result.LastInsertId()
}

func UpdateUser(db *sql.DB, user User) (int64, error) {
	result, err := db.Exec("UPDATE gateway_user SET Password = ?, PublicKey = ?, secret_key = ?, Salt = ? WHERE Id = ?", user.Password, user.PublicKey, user.SecretKey, user.Salt, user.Id)
	if err != nil {
		zap.L().Error("Error updating user", zap.Error(err))
		return 0, err
	}
	return result.RowsAffected()
}

func RemoveUser(db *sql.DB, id int) (int64, error) {
	result, err := db.Exec("DELETE FROM gateway_user WHERE Id = ?", id)
	if err != nil {
		zap.L().Error("Error deleting user", zap.Error(err))
		return 0, err
	}
	return result.RowsAffected()
}

func GetUsers(db *sql.DB) ([]User, error) {

	rows, err := db.Query("SELECT * FROM gateway_user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Salt, &user.Password, &user.PublicKey, &user.SecretKey); err != nil {
			zap.L().Error("Error getting all users", zap.Error(err))
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetUserByPublicKey(db *sql.DB, publicKey string) (User, error) {
	var user User
	rows := db.QueryRow("SELECT * FROM gateway_user WHERE Public_Key = ?", publicKey)
	if err := rows.Scan(&user.Id, &user.Salt, &user.Password, &user.PublicKey, &user.SecretKey); err != nil {
		zap.L().Error("Error getting user by public key", zap.Error(err))
		return User{}, err
	}
	return user, nil
}

func GetUserByPassword(db *sql.DB, password string) (User, error) {
	var user User
	rows := db.QueryRow("SELECT * FROM gateway_user WHERE Password = ?", password)
	if err := rows.Scan(&user.Id, &user.Salt, &user.Password, &user.PublicKey, &user.SecretKey); err != nil {
		zap.L().Error("Error getting user by password", zap.Error(err))
		return User{}, err
	}
	return user, nil
}

func GetUserByID(db *sql.DB, id int) (User, error) {
	var user User
	rows := db.QueryRow("SELECT * FROM gateway_user WHERE Id = ?", id)
	if err := rows.Scan(&user.Id, &user.Salt, &user.Password, &user.PublicKey, &user.SecretKey); err != nil {
		zap.L().Error("Error getting user by Id", zap.Error(err))
		return User{}, err
	}
	return user, nil
}
