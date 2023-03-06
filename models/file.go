package models

import (
	"fmt"

	uuid "github.com/nu7hatch/gouuid"
	"gorm.io/gorm"
)

type Model interface {
	GetId() string
	SetId(id string)
}

type File struct {
	Id          string `json:"id"`
	Key         string `json:"key"`
	FileName    string `json:"fileName"`
	Bucket      string `json:"bucket"`
	ContentType string `json:"content-type"`
	Size        int64  `json:"size"`
}

func (*File) NewFile(key string, filename string, bucket string, contentType string, size int64) *File {

	db, err := dbConnection()
	// Auto Migrate
	db.AutoMigrate(&File{})
	// Set table options
	db.Set("gorm:table_options", "ENGINE=Distributed(cluster, default, hits)").AutoMigrate(&File{})

	u, err2 := uuid.NewV4()

	file := &File{
		Id:          u.String(),
		Key:         key,
		FileName:    filename,
		Bucket:      bucket,
		ContentType: contentType,
		Size:        size,
	}

	// Insert
	db.Create(&file)

	if err2 != nil {
		fmt.Printf("could not insert row: %v", err)
		panic(err)
	}

	return file
}

func GetAllFiles() *gorm.DB {
	var files []File

	db, err := dbConnection()

	if err != nil {
		fmt.Printf("could not insert row: %v", err)
		panic(err)
	}

	result := db.Find(&files)

	return result
}

func (u *File) GetId() string {
	return u.FileName
}
