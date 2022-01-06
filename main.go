package main

import (
	"database/sql"
	"fmt"

	"log"
	"net/http"

	"github.com/Thauan/golang-api-upload-s3-example/controllers"
	"github.com/Thauan/golang-api-upload-s3-example/handlers"
	"github.com/Thauan/golang-api-upload-s3-example/middlewares"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	s3session *s3.S3
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

func init() {
	handlers.LoadEnv()
	awsAccessKeyID := handlers.GetEnvWithKey("AWS_ACCESS_KEY_ID")
	awsSecretAccess := handlers.GetEnvWithKey("AWS_SECRET_ACCESS_KEY")
	MyRegion := handlers.GetEnvWithKey("AWS_S3_REGION")
	DatabasePort := handlers.GetEnvWithKey("DATABASE_PORT")
	DatabaseHost := handlers.GetEnvWithKey("DATABASE_HOST")
	DatabaseTable := handlers.GetEnvWithKey("DATABASE_TABLE")
	DatabaseUser := handlers.GetEnvWithKey("DATABASE_USER")
	DatabasePassword := handlers.GetEnvWithKey("DATABASE_PASSWORD")
	sslMode := handlers.GetEnvWithKey("SSL_MODE")

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

func main() {

	Port := handlers.GetEnvWithKey("PORT")
	fmt.Println("Running in http://localhost:" + Port)

	router := mux.NewRouter()

	router.HandleFunc("/files", controllers.GetFiles(s3session)).Methods("GET")
	router.Handle("/upload", middlewares.IsAuthorized(controllers.UploadFiles(s3session))).Methods("POST")
	router.Handle("/generate/thumb", middlewares.IsAuthorized(controllers.GenerateThumbVideo(s3session))).Methods("POST")

	log.Fatal(http.ListenAndServe(":"+Port, router))
}
