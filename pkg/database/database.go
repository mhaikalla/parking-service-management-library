package models

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
)

// IServerDB ..
type IServerDB struct {
	DB *gorm.DB
}

// NewDBConnection ..
func NewDBConnection(config map[string]map[string]interface{}) IServerDB {
	if c, f := config["gorm"]; f {
		Dbdriver := c["dialect"].(string)
		DBURL := c["connectionstring"].(string)
		connection, err := gorm.Open(Dbdriver, DBURL)
		//connection.LogMode(true)
		if err != nil {
			errMsg := fmt.Sprintf("Cannot connect to %s database", Dbdriver)
			panic(errors.New(errMsg))
		}
		return IServerDB{DB: connection}
	}
	panic(errors.New("config not found"))
}
