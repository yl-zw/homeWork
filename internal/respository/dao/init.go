package dao

import (
	"database/sql"
	"fmt"
	"webbook/config"
)

func InitTables() error {
	//err := db.AutoMigrate(&user.User{})
	//if err != nil {
	//	return err
	//}
	//err = db.AutoMigrate(&user.Profile{})
	//if err != nil {
	//	return err
	//}
	conn, err := sql.Open("mysql", config.Con.DB.Dns(true))
	if err != nil {
		return err
	}
	defer func(conn *sql.DB) {
		err = conn.Close()
		if err != nil {
			return
		}
	}(conn)
	_, err = conn.Exec(databases)
	if err != nil {
		return err
	}
	res, err := conn.Exec(initSqlUser)
	if err != nil {
		fmt.Println(res)
		return err
	}
	_, err = conn.Exec(initProfile)
	if err != nil {
		return err
	}
	return nil
}
