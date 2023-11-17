package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Thauan/golang-api-upload-s3-example/handlers"
	"github.com/Thauan/golang-api-upload-s3-example/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetFiles(session *s3.S3) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files := handlers.GetObjectsStorage(session)

		data, _ := json.Marshal(files)

		w.Write(data)
	}
}

func GetExampleFunc(session *s3.S3) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("funcionou")
	}
}

func UploadFiles(session *s3.S3) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var File models.File

		var files []models.File

		file, _ := r.MultipartReader()

		fmt.Println(file, "file")

		for {
			part, err := file.NextPart()
			if err == io.EOF {
				break
			}
			fileBytes, err2 := ioutil.ReadAll(part)

			if err2 != nil {
				data, _ := json.Marshal(fmt.Sprintf("failed to upload file %v", err2))
				w.WriteHeader(http.StatusBadGateway)
				w.Write(data)
				return
			}

			filename := "upload-*" + string(filepath.Ext(part.FileName()))

			tempFile, err3 := ioutil.TempFile(os.TempDir(), filename)

			if err3 != nil {
				data, _ := json.Marshal(fmt.Sprintf("failed to upload file %v", err3))
				w.WriteHeader(http.StatusBadGateway)
				w.Write(data)
				return
			}

			if part != nil {
				part.Close()
			}

			defer tempFile.Close()

			defer handlers.RemoveTempFile(tempFile)

			tempFile.Write(fileBytes)

			fmt.Println("Done upload temp file")

			resp, size := handlers.MultipartUploadObject(session, tempFile.Name())

			fmt.Println(resp)

			new := File.NewFile(*resp.Key, *resp.Location, *resp.Bucket, part.Header.Get("Content-Type"), size)

			files = append(files, *new)
		}

		data, _ := json.Marshal(files)

		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}
}

func GenerateThumbVideo(session *s3.S3) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var File models.File

		frame := r.FormValue("frame")
		file, handler, err := r.FormFile("file")

		if err != nil {
			data, _ := json.Marshal(fmt.Sprintf("failed to upload file %v", err))
			w.WriteHeader(http.StatusBadGateway)
			w.Write(data)
			return
		}

		file.Close()

		filename := "upload-*" + string(filepath.Ext(handler.Filename))

		tempFile, err2 := os.CreateTemp(os.TempDir(), filename)

		if err2 != nil {
			data, _ := json.Marshal(fmt.Sprintf("failed to upload file %v", err2))
			w.WriteHeader(http.StatusBadGateway)
			w.Write(data)
			return
		}

		defer tempFile.Close()

		// fileBytes, err3 := io.ReadAll(file)

		// if err3 != nil {
		// 	data, _ := json.Marshal(fmt.Sprintf("failed to upload file %v", err3))
		// 	w.WriteHeader(http.StatusBadGateway)
		// 	w.Write(data)
		// 	return
		// }

		frameInt, _ := strconv.Atoi(frame)

		thumb, err4 := handlers.TakeFrame(tempFile.Name(), frameInt)

		var bytesWrite bytes.Buffer

		bytesWrite.ReadFrom(thumb)

		tempFile.Write([]byte(bytesWrite.Bytes()))

		fmt.Println("Done upload temp file")

		if err4 != nil {
			data, _ := json.Marshal(fmt.Sprintf("failed to generate thumb file %v", err4.Error()))
			w.WriteHeader(http.StatusBadGateway)
			w.Write(data)
			return
		}

		// resp, size := handlers.MultipartUploadObject(session, tempFile.Name())
		resp, size := handlers.MultipartUploadObject(session, tempFile.Name())

		params := &s3.GetObjectInput{
			Bucket: aws.String(handlers.GetEnvWithKey("AWS_S3_BUCKET")),
			Key:    aws.String(*resp.Key),
		}

		req, _ := session.GetObjectRequest(params)
		url, err5 := req.Presign(15 * time.Minute)

		if err5 != nil {
			data, _ := json.Marshal(fmt.Sprintf("error in presign: %v", err5.Error()))
			w.WriteHeader(http.StatusBadGateway)
			w.Write(data)
			return
		}

		db := File.NewFile(*resp.Key, url, *resp.Bucket, handler.Header.Get("Content-Type"), size)

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
