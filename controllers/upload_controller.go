package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Thauan/golang-api-upload-s3-example/handlers"
	"github.com/Thauan/golang-api-upload-s3-example/models"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetFiles(session *s3.S3) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files := handlers.GetObjectsStorage(session)

		data, _ := json.Marshal(files)

		w.Write(data)
	}
}

func UploadFiles(session *s3.S3) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var File models.File

		file, handler, err := r.FormFile("file")

		if err != nil {
			data, _ := json.Marshal(fmt.Sprintf("failed to upload file %v", err))
			w.WriteHeader(http.StatusBadGateway)
			w.Write(data)
			return
		}

		filename := "upload-*" + string(filepath.Ext(handler.Filename))

		tempFile, err2 := ioutil.TempFile(os.TempDir(), filename)

		if err2 != nil {
			data, _ := json.Marshal(fmt.Sprintf("failed to upload file %v", err2))
			w.WriteHeader(http.StatusBadGateway)
			w.Write(data)
			return
		}

		defer tempFile.Close()

		defer handlers.RemoveTempFile(tempFile)

		fileBytes, err3 := ioutil.ReadAll(file)

		if err3 != nil {
			data, _ := json.Marshal(fmt.Sprintf("failed to upload file %v", err3))
			w.WriteHeader(http.StatusBadGateway)
			w.Write(data)
			return
		}

		tempFile.Write(fileBytes)

		fmt.Println("Done upload temp file")

		resp, size := handlers.MultipartUploadObject(session, tempFile.Name())

		db := File.NewFile(*resp.Key, *resp.Location, *resp.Bucket, handler.Header.Get("Content-Type"), size)

		data, _ := json.Marshal(db)

		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}
}

func GenerateThumbVideo(session *s3.S3) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var File models.File

		file, handler, err := r.FormFile("file")

		if err != nil {
			data, _ := json.Marshal(fmt.Sprintf("failed to upload file %v", err))
			w.WriteHeader(http.StatusBadGateway)
			w.Write(data)
			return
		}

		filename := "upload-*" + string(filepath.Ext(handler.Filename))

		tempFile, err2 := ioutil.TempFile(os.TempDir(), filename)

		if err2 != nil {
			data, _ := json.Marshal(fmt.Sprintf("failed to upload file %v", err2))
			w.WriteHeader(http.StatusBadGateway)
			w.Write(data)
			return
		}

		defer tempFile.Close()

		// defer handlers.RemoveTempFile(tempFile)

		fileBytes, err3 := ioutil.ReadAll(file)

		if err3 != nil {
			data, _ := json.Marshal(fmt.Sprintf("failed to upload file %v", err3))
			w.WriteHeader(http.StatusBadGateway)
			w.Write(data)
			return
		}

		thumb, err4 := handlers.FirstFrame(tempFile.Name(), fileBytes)

		var bytesWrite bytes.Buffer

		bytesWrite.ReadFrom(thumb)

		tempFile.Write([]byte(bytesWrite.String()))

		fmt.Println("Done upload temp file")

		if err4 != nil {
			data, _ := json.Marshal(fmt.Sprintf("failed to generate thumb file %v", err4.Error()))
			w.WriteHeader(http.StatusBadGateway)
			w.Write(data)
			return
		}

		resp, size := handlers.MultipartUploadObject(session, tempFile.Name())

		db := File.NewFile(*resp.Key, *resp.Location, *resp.Bucket, handler.Header.Get("Content-Type"), size)

		data, _ := json.Marshal(db)

		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}
}

func ConvertVideoToMP4(session *s3.S3) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var File models.File

		file, handler, err := r.FormFile("file")

		if err != nil {
			data, _ := json.Marshal(fmt.Sprintf("failed to upload file %v", err))
			w.WriteHeader(http.StatusBadGateway)
			w.Write(data)
			return
		}

		filename := "upload-*" + string(filepath.Ext(handler.Filename))

		tempFile, err2 := ioutil.TempFile(os.TempDir(), filename)

		if err2 != nil {
			data, _ := json.Marshal(fmt.Sprintf("failed to upload file %v", err2))
			w.WriteHeader(http.StatusBadGateway)
			w.Write(data)
			return
		}

		defer tempFile.Close()

		// defer handlers.RemoveTempFile(tempFile)

		fileBytes, err3 := ioutil.ReadAll(file)

		if err3 != nil {
			data, _ := json.Marshal(fmt.Sprintf("failed to upload file %v", err3))
			w.WriteHeader(http.StatusBadGateway)
			w.Write(data)
			return
		}

		thumb, err4 := handlers.VideoToMP4(tempFile.Name(), fileBytes)

		var bytesWrite bytes.Buffer

		bytesWrite.ReadFrom(thumb)

		tempFile.Write([]byte(bytesWrite.String()))

		fmt.Println("Done upload temp file")

		if err4 != nil {
			data, _ := json.Marshal(fmt.Sprintf("failed to generate thumb file %v", err4.Error()))
			w.WriteHeader(http.StatusBadGateway)
			w.Write(data)
			return
		}

		resp, size := handlers.MultipartUploadObject(session, tempFile.Name())

		db := File.NewFile(*resp.Key, *resp.Location, *resp.Bucket, handler.Header.Get("Content-Type"), size)

		data, _ := json.Marshal(db)

		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}
}
