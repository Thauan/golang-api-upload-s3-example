package models

import (
	"fmt"

	"github.com/Thauan/golang-api-upload-s3-example/utils"
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
	FileUrl     string `json:"fileUrl"`
	FileFolder  string `json:"fileFolder"`
	Bucket      string `json:"bucket"`
	ContentType string `json:"content-type"`
	Size        int64  `json:"size"`
}

func (*File) NewFile(key string, url string, bucket string, contentType string, size int64) *File {

	db, err := dbConnection()
	// Auto Migrate
	db.AutoMigrate(&File{})
	// Set table options
	db.Set("gorm:table_options", "ENGINE=Distributed(cluster, default, hits)").AutoMigrate(&File{})

	if err != nil {
		fmt.Printf("could not connect: %v", err)
	}

	u, err2 := uuid.NewV4()

	dir, f := utils.GetFileWithDir(key)

	file := &File{
		Id:          u.String(),
		Key:         key,
		FileName:    f,
		FileFolder:  dir,
		FileUrl:     url,
		Bucket:      bucket,
		ContentType: contentType,
		Size:        size,
	}

	// Insert
	db.Create(&file)

	if err2 != nil {
		fmt.Printf("could not insert row: %v", err2)
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

func GetNextID() (int64, error) {
	// Crie uma instância do objeto DB do GORM
	db, err := dbConnection()

	if err != nil {
		return 0, err
	}

	// Use o método AutoMigrate() do GORM para criar a tabela, se ainda não existir
	db.AutoMigrate(&File{})
	// Set table options
	db.Set("gorm:table_options", "ENGINE=Distributed(cluster, default, hits)").AutoMigrate(&File{})

	var maxID int64

	// Use o método Select() do GORM para selecionar o maior ID
	if err := db.Table("files").Select("MAX(id)").Scan(&maxID).Error; err != nil {
		return 0, err
	}

	// Retorna o próximo ID
	return maxID + 1, nil
}
