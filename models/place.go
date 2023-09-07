package models

import (
	"fmt"

	uuid "github.com/nu7hatch/gouuid"
	"gorm.io/gorm"
)

type Place struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (*File) NewPlace(key string, name string) *Place {

	db, err := dbConnection()
	// Auto Migrate
	db.AutoMigrate(&Place{})
	// Set table options
	db.Set("gorm:table_options", "ENGINE=Distributed(cluster, default, hits)").AutoMigrate(&File{})

	u, err2 := uuid.NewV4()

	place := &Place{
		Id:   u.String(),
		Name: name,
	}

	// Insert
	db.Create(&place)

	if err2 != nil {
		fmt.Printf("could not insert row: %v", err)
		panic(err)
	}

	return place
}

func GetAllPlaces() *gorm.DB {
	var places []Place

	db, err := dbConnection()

	if err != nil {
		fmt.Printf("could not insert row: %v", err)
		panic(err)
	}

	result := db.Find(&places)

	return result
}
