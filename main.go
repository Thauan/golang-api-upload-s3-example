package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"golang-api-upload-s3-example.com/uploader/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	s3session *s3.S3
)

const (
	FILE      = "files/teste.txt"
	PART_SIZE = 6_000_000 // Has to be 5_000_000 minimim
	RETRIES   = 2
)

var AccessKeyID string
var SecretAccessKey string
var MyRegion string
var DatabaseUrl string
var DatabaseHost string
var DatabasePort string
var DatabaseTable string
var DatabaseUser string
var DatabasePassword string
var Port string
var sslMode string

func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
}

func init() {
	LoadEnv()
	awsAccessKeyID := GetEnvWithKey("AWS_ACCESS_KEY_ID")
	awsSecretAccess := GetEnvWithKey("AWS_SECRET_ACCESS_KEY")
	MyRegion := GetEnvWithKey("AWS_S3_REGION")
	DatabasePort := GetEnvWithKey("DATABASE_PORT")
	DatabaseHost := GetEnvWithKey("DATABASE_HOST")
	DatabaseTable := GetEnvWithKey("DATABASE_TABLE")
	DatabaseUser := GetEnvWithKey("DATABASE_USER")
	DatabasePassword := GetEnvWithKey("DATABASE_PASSWORD")
	sslMode := GetEnvWithKey("SSL_MODE")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		DatabaseHost, DatabasePort, DatabaseUser, DatabasePassword, DatabaseTable, sslMode)

	fmt.Println(psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		fmt.Println("You not connect to your database")
		panic(err)
	}

	fmt.Println("You connected to your database: " + DatabaseTable)

	awsConfig := &aws.Config{
		Region:      aws.String(MyRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccess, ""),
	}

	awsConfig.WithLogLevel(aws.LogDebug)

	s3session = s3.New(session.Must(session.NewSession(awsConfig)))
}

func MultipartUploadObject(filename string) (result *s3.CompleteMultipartUploadOutput) {
	// Open a file.
	file, err := os.Open(filename)
	defer file.Close()

	// Get file size
	stats, _ := file.Stat()
	size := stats.Size()

	// put file in byteArray
	buffer := make([]byte, size)
	file.Read(buffer)

	// Create MultipartUpload object
	expiryDate := time.Now().AddDate(0, 0, 1)

	createdResp, err := s3session.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket:  aws.String(GetEnvWithKey("AWS_S3_BUCKET")),
		Key:     aws.String(file.Name()),
		Expires: &expiryDate,
	})

	var start, currentSize int
	var remaining = int(size)
	var partNum = 1
	var completedParts []*s3.CompletedPart

	// Loop till remaining upload size is 0
	for start = 0; remaining != 0; start += PART_SIZE {
		if remaining < PART_SIZE {
			currentSize = remaining
		} else {
			currentSize = PART_SIZE
		}

		completed, err := Upload(createdResp, buffer[start:start+currentSize], partNum)
		// If upload function failed (meaning it retried acoording to RETRIES)
		if err != nil {
			_, err = s3session.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
				Bucket:   createdResp.Bucket,
				Key:      createdResp.Key,
				UploadId: createdResp.UploadId,
			})
			if err != nil {
				// god speed
				fmt.Println(err)
				return
			}
		}

		// Detract the current part size from remaining
		remaining -= currentSize
		fmt.Printf("Part %v complete, %v btyes remaining\n", partNum, remaining)

		// Add the completed part to our list
		completedParts = append(completedParts, completed)
		partNum++

	}

	// All the parts are uploaded, completing the upload
	result, err = s3session.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:   createdResp.Bucket,
		Key:      createdResp.Key,
		UploadId: createdResp.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result.String())
	}

	return result
}

func listObjects() (resp *s3.ListObjectsV2Output) {
	resp, err := s3session.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String("golang-s3"),
	})

	if err != nil {
		panic(err)
	}

	return resp
}

func main() {

	Port := GetEnvWithKey("PORT")
	fmt.Println("Running in http://localhost:" + Port)

	// ROUTER CONFIGURATION
	router := mux.NewRouter()
	router.HandleFunc("/files", GetFiles).Methods("GET")
	router.HandleFunc("/upload", UploadFiles).Methods("POST")

	log.Fatal(http.ListenAndServe(":"+Port, router))
}

func GetFiles(w http.ResponseWriter, r *http.Request) {

	files := listObjects()

	fmt.Printf("Resultado:%s", files)

	data, _ := json.Marshal(files)

	w.Write(data)

}

func UploadFiles(w http.ResponseWriter, r *http.Request) {

	var File models.File

	fmt.Println(File)

	_, handler, err := r.FormFile("file")

	if err != nil {
		fmt.Println(handler)
		data, _ := json.Marshal(fmt.Sprintf("failed to open file %q, %v", handler.Filename, err))
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}

	resp := MultipartUploadObject(FILE)

	fmt.Println(handler.Filename)

	data, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func uploadObject(filename string) (resp *s3.PutObjectOutput) {
	bucket := GetEnvWithKey("AWS_S3_BUCKET")

	file, err := os.Open(filename)
	defer file.Close()
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()

	if err != nil {
		fmt.Println(err)
		return
	}

	buffer := make([]byte, size)
	_, _ = file.Read(buffer)

	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	params := &s3.PutObjectInput{
		Body:        fileBytes,
		Bucket:      aws.String(bucket),
		Key:         aws.String(strings.Split(filename, "/")[1]),
		ContentType: aws.String(fileType),
		ACL:         aws.String(s3.BucketCannedACLPublicRead),
	}

	resp, _ = s3session.PutObject(params)

	return resp
}

// Uploads the fileBytes bytearray a MultiPart upload
func Upload(resp *s3.CreateMultipartUploadOutput, fileBytes []byte, partNum int) (completedPart *s3.CompletedPart, err error) {
	var try int
	for try <= RETRIES {
		uploadResp, err := s3session.UploadPart(&s3.UploadPartInput{
			Body:          bytes.NewReader(fileBytes),
			Bucket:        resp.Bucket,
			Key:           resp.Key,
			PartNumber:    aws.Int64(int64(partNum)),
			UploadId:      resp.UploadId,
			ContentLength: aws.Int64(int64(len(fileBytes))),
		})
		// Upload failed
		if err != nil {
			fmt.Println(err)
			// Max retries reached! Quitting
			if try == RETRIES {
				return nil, err
			} else {
				// Retrying
				try++
			}
		} else {
			// Upload is done!
			return &s3.CompletedPart{
				ETag:       uploadResp.ETag,
				PartNumber: aws.Int64(int64(partNum)),
			}, nil
		}
	}

	return nil, nil
}
