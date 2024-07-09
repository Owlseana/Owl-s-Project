// database/db.go
package dao

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	mymongo "github.com/quan-xie/tuba/database/mongo"
	"github.com/quan-xie/tuba/util/xtime"
	xmongo "go.mongodb.org/mongo-driver/mongo"
)

var DB *sql.DB
var client *xmongo.Client

func Initmysql() {
	var err error
	DB, err = sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/mydatabase?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	// Test the connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("failed to ping database:", err)
	}
}

func InitMongo() {
	cfg := &mymongo.Config{
		Addrs:         []string{"localhost:27017"},
		Username:      "",
		Password:      "",
		MaxPool:       2,
		ReplicaSet:    "",
		MaxIdletime:   xtime.Duration(time.Second * 10),
		SocketTimeout: xtime.Duration(500 * time.Millisecond),
		ConnTimeout:   xtime.Duration(500 * time.Millisecond),
	}
	client = mymongo.NewMongo(cfg)
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connected to MongoDB!")
}
