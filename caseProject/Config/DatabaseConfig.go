package Config

import (
	"caseProject/DataModel"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func InitDatabaseConfig() {

	var err error
	Database, err = gorm.Open(postgres.Open("host=localhost user=postgres password=case dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Istanbul"), &gorm.Config{})

	if err != nil {
		fmt.Println("Gorm DB Open Error", err)
		panic(err)
	}

	err = Database.AutoMigrate(
		&DataModel.Message{},
		&DataModel.MessageOutbox{})

}
