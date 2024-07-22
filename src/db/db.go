package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/kylerequez/go-upload-example/src/utils"
)

var DB *sql.DB

func ConnectDB() error {
	dsn, err := GenerateDSN()
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func PingDB() error {
	return DB.Ping()
}

func GenerateDSN() (string, error) {
	errFmt := "%s is empty"

	DSN := "DB_DSN"
	dsn := utils.GetEnv(DSN)
	if dsn == "" {
		return "", fmt.Errorf(errFmt, DSN)
	}

	HOST := "DB_HOST"
	host := utils.GetEnv(HOST)
	if host == "" {
		return "", fmt.Errorf(errFmt, HOST)
	}

	// PORT := "DB_PORT"
	// port := utils.GetEnv(PORT)
	// if host == "" {
	// 	return "", fmt.Errorf(errFmt, PORT)
	// }

	USERNAME := "DB_USERNAME"
	username := utils.GetEnv(USERNAME)
	if host == "" {
		return "", fmt.Errorf(errFmt, USERNAME)
	}

	PASSWORD := "DB_PASSWORD"
	password := utils.GetEnv(PASSWORD)
	if host == "" {
		return "", fmt.Errorf(errFmt, PASSWORD)
	}

	NAME := "DB_NAME"
	name := utils.GetEnv(NAME)
	if host == "" {
		return "", fmt.Errorf(errFmt, NAME)
	}

	return fmt.Sprintf(dsn, username, password, host, name), nil
}

// type Database struct {
// 	Conn     any
// 	Dsn      string
// 	Server   string
// 	Port     string
// 	Username string
// 	Password string
// }
//
// func NewDatabase() *Database {
// 	return &Database{}
// }
//
// func (db *Database) ConnectDB() error {
// 	errorString := "%s is empty"
//
// 	if db.Server == "" {
// 		return fmt.Errorf(errorString, "server")
// 	}
//
// 	if db.Port == "" {
// 		return fmt.Errorf(errorString, "port")
// 	}
//
// 	if db.Username == "" {
// 		return fmt.Errorf(errorString, "username")
// 	}
//
// 	if db.Password == "" {
// 		return fmt.Errorf(errorString, "password")
// 	}
//
// 	if db.Dsn != "" {
// 		dsn := utils.GetEnv("DSN")
// 		if dsn == "" {
// 			return fmt.Errorf("%s is empty", "dsn")
// 		}
//
// 		db.Dsn = fmt.Sprintf(
// 			dsn,
// 			db.Server,
// 			db.Port,
// 			db.Username,
// 			db.Password,
// 		)
// 	}
//
// 	return nil
// }
//
// func (db *Database) SetServer(server string) *Database {
// 	db.Server = server
// 	return db
// }
//
// func (db *Database) SetPort(port string) *Database {
// 	db.Port = port
// 	return db
// }
//
// func (db *Database) SetUsername(username string) *Database {
// 	db.Username = username
// 	return db
// }
//
// func (db *Database) SetPassword(password string) *Database {
// 	db.Password = password
// 	return db
// }
