package models

import (
	"database/sql"
)

type Model interface {
	GetId() string
	SetId(id string)
}

type modelImpl struct {
	id string
}

type File struct {
	modelImpl   `json:"id"`
	Key         string `json:"key"`
	FileName    string `json:"fileName"`
	Bucket      string `json:"bucket"`
	ContentType string `json:"content-type"`
	Size        int64  `json:"size"`
}

func (m *modelImpl) SetId(id string) {
	m.id = id
}

func (*File) NewFile(key string, filename string, bucket string, contentType string, size int64) *File {

	// var DB *sql.DB

	f := &File{
		Key:         key,
		FileName:    filename,
		Bucket:      bucket,
		ContentType: contentType,
		Size:        size,
	}

	// result, err := DB.Exec("INSERT INTO files (id, filename, content_type, file_size, bucket) VALUES ($1, $2, $3, $4. $5)", f.id, f.FileName, f.ContentType, f.Size, f.Bucket)

	// if err != nil {
	// 	fmt.Printf("could not insert row: %v", err)
	// 	panic(err)
	// }

	f.SetId(filename)

	// rowsAffected, err := result.RowsAffected()

	// we can log how many rows were inserted
	// fmt.Println("inserted", rowsAffected, "rows")

	return f
}

// AllBooks returns a slice of all books in the books table.
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
