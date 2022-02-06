package models

import (
	"database/sql"
	"fmt"

	"github.com/Thauan/golang-api-upload-s3-example/handlers"
	uuid "github.com/nu7hatch/gouuid"
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

	handlers.LoadEnv()
	DatabasePort := handlers.GetEnvWithKey("DATABASE_PORT")
	DatabaseHost := handlers.GetEnvWithKey("DATABASE_HOST")
	DatabaseTable := handlers.GetEnvWithKey("DATABASE_TABLE")
	DatabaseUser := handlers.GetEnvWithKey("DATABASE_USER")
	DatabasePassword := handlers.GetEnvWithKey("DATABASE_PASSWORD")
	sslMode := handlers.GetEnvWithKey("SSL_MODE")

	u, err2 := uuid.NewV4()

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		DatabaseHost, DatabasePort, DatabaseUser, DatabasePassword, DatabaseTable, sslMode)

	fmt.Println(psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)

	if err2 != nil {
		fmt.Println(err2)
	}

	f := &File{
		Id:          u.String(),
		Key:         key,
		FileName:    filename,
		Bucket:      bucket,
		ContentType: contentType,
		Size:        size,
	}

	result, err := db.Prepare("INSERT INTO medias (id, path, filename, bucket, type, size) VALUES ($1, $2, $3, $4, $5, $6)")

	if err != nil {
		fmt.Printf("could not insert row: %v", err)
		panic(err)
	}

	_, err3 := result.Exec(f.Id, f.Key, f.FileName, f.Bucket, f.ContentType, f.Size)

	if err3 != nil {
		fmt.Printf("could not insert row: %v", err)
		panic(err)
	}

	defer db.Close()

	return f
}

func getAllFiles() ([]File, error) {
	// Create an exported global variable to hold the database connection pool.
	var DB *sql.DB
	// Note that we are calling Query() on the global variable.
	rows, err := DB.Query("SELECT * FROM files")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File

	for rows.Next() {
		var file File

		err := rows.Scan(&file.Key, &file.FileName, &file.Bucket)

		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func (u *File) GetId() string {
	return u.FileName
}
