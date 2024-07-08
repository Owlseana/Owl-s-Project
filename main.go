package main

import (
	"MyFirstProgram/dao"
	"MyFirstProgram/router"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dao.Initmysql()
	// dao.InitMongo()
	router.Initrouter()
}
