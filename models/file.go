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
	modelImpl
	Key      string //image_content
	FileName string //test.jpg
	Content  []byte //[]byte
	path     string
	bucket   string
}

func (m *modelImpl) SetId(id string) {
	m.id = id
}

func NewFile(key string, filename string, content []byte, path string, bucket string, name string) *File {
	u := &File{
		Key:      key,
		FileName: filename,
		Content:  content,
		path:     path,
		bucket:   bucket,
	}
	u.SetId(name)
	return u
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

		err := rows.Scan(&file.Key, &file.path, &file.FileName, &file.Content, &file.bucket)
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
