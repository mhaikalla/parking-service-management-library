package migration

import (
	"log"

	modelsDB "github.com/mhaikalla/parking-service-management-library/pkg/database"
)

func RunMigration(dbconnection modelsDB.IServerDB) {
	db := dbconnection.DB
	if db.Error != nil {
		log.Fatalln(db.Error.Error())
	}

	// tableName := mConsent.ConsentStatus{}.ConsentStatusTableName()
	// if exist := db.HasTable(tableName); !exist {
	// 	fmt.Println("migrate table " + tableName)
	// 	err := db.CreateTable(&mConsent.ConsentStatus{})
	// 	if err == nil {
	// 		fmt.Println("success migrate table " + tableName)
	// 	}
	// }

}
