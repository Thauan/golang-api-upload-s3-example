package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	// "encoding/json"
	// "github.com/bitly/go-simplejson"
	// "github.com/aws/aws-sdk-go/aws/awsutil" 
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
)

var (s3session *s3.S3)
var AccessKeyID string
var SecretAccessKey string
var MyRegion string


//GetEnvWithKey : get env value
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

	awsConfig := &aws.Config{
		Region: aws.String(MyRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccess, ""),
	}

	awsConfig.WithLogLevel(aws.LogDebug)

	s3session = s3.New(session.Must(session.NewSession(awsConfig)))
}

func uploadObject(filename string) ( resp *s3.PutObjectOutput) {
	bucket := GetEnvWithKey("AWS_S3_BUCKET")

	file, err := os.Open(filename)
	defer file.Close()
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()

	buffer := make([]byte, size)
	_, _ = file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	if err != nil {
		panic(err)
	}

	fmt.Println("Uploading:", filename)

	params := &s3.PutObjectInput{
		Body: fileBytes,
		Bucket: aws.String(bucket),
		Key: aws.String(strings.Split(filename, "/")[1]),
		ContentType: aws.String(fileType),
		ACL: aws.String(s3.BucketCannedACLPublicRead),
	}

	resp, err = s3session.PutObject(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	return resp
}

func listObjects() (resp *s3.ListObjectsV2Output) {
	resp, err := s3session.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String("golang-s3"),
	})

	if err != nil {
		panic(err)
	}

	// fmt.Printf("response %s", awsutil.StringValue(resp)) 

	return resp
}

func main() {
	fmt.Printf("Running S3 uploader in golang");
	
	// ROUTER CONFIGURATION
	router := mux.NewRouter()
	router.HandleFunc("/files", GetFiles).Methods("GET")
	router.HandleFunc("/upload", UploadFiles).Methods("POST")

    log.Fatal(http.ListenAndServe(":8000", router))
}

func GetFiles(w http.ResponseWriter, r *http.Request) {
	// função de retorno dos files
	files := listObjects() 
	// (files *s3.ListObjectsV2Output)

	// reader := strings.NewReader(files)
	
	// dec := json.NewDecoder(reader)
	
	fmt.Printf("Resultado:%s", files)
	// return json.NewJson(files, []byte)

}

func UploadFiles(w http.ResponseWriter, r *http.Request) {
	// função de retorno dos files
	uploadObject("files/teste.txt")
}