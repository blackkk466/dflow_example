package main

import (
	"log"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"fmt"
)

var db,dberr = gorm.Open("mysql", "root:123456@/database1?charset=utf8&parseTime=True&loc=Local")


func main() {

	connectToDatabase()
	migrateModel()
	defer db.Close()

	r := gin.Default()

	r.GET("/webhook",verificationEndPoint)
	r.POST("/webhook",messagesEndPoint)
	r.POST("/apiai",APIAIPopulationEndpoint)

	log.Fatal(r.Run(":8888"))

}

func connectToDatabase() {
	if dberr != nil {
		fmt.Printf("it dieded")
		log.Fatal(dberr)
	}
	db.SingularTable(true)
}

func migrateModel(){
	DB := db

	//DB.DropTableIfExists(&User{},&Note{})

	DB.AutoMigrate(
		&Note{},
		&Product{},
		&Session{},)
}