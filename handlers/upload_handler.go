package handlers

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	PART_SIZE = 6_000_000
	RETRIES   = 2
)

func GetObjectsStorage(session *s3.S3) (resp *s3.ListObjectsV2Output) {
	resp, err := session.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String("golang-s3"),
	})

	if err != nil {
		panic(err)
	}

	return resp
}

func MultipartUploadObject(session *s3.S3, filename string) (result *s3.CompleteMultipartUploadOutput, completeSize int64) {

	file, err := os.Open(filename)
	defer file.Close()

	stats, _ := file.Stat()
	size := stats.Size()

	buffer := make([]byte, size)
	file.Read(buffer)

	expiryDate := time.Now().AddDate(0, 0, 1)

	createdResp, err := session.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket:  aws.String(GetEnvWithKey("AWS_S3_BUCKET")),
		Key:     aws.String(file.Name()),
		Expires: &expiryDate,
	})

	var start, currentSize int
	var remaining = int(size)
	var partNum = 1
	var completedParts []*s3.CompletedPart

	for start = 0; remaining != 0; start += PART_SIZE {
		if remaining < PART_SIZE {
			currentSize = remaining
		} else {
			currentSize = PART_SIZE
		}

		completed, err := Upload(session, createdResp, buffer[start:start+currentSize], partNum)

		if err != nil {
			_, err = session.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
				Bucket:   createdResp.Bucket,
				Key:      createdResp.Key,
				UploadId: createdResp.UploadId,
			})
			if err != nil {
				// fmt.Println(err)
				return
			}
		}

		remaining -= currentSize
		// fmt.Printf("Part %v complete, %v bytes remaining\n", partNum, remaining)

		completedParts = append(completedParts, completed)
		partNum++

	}

	result, err = session.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:   createdResp.Bucket,
		Key:      createdResp.Key,
		UploadId: createdResp.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})

	completeSize = size

	if err != nil {
		// fmt.Println(err)
	} else {
		// fmt.Println(result.String())
	}

	return result, completeSize
}

// Uploads the fileBytes bytearray a MultiPart upload
func Upload(session *s3.S3, resp *s3.CreateMultipartUploadOutput, fileBytes []byte, partNum int) (completedPart *s3.CompletedPart, err error) {
	var try int
	for try <= RETRIES {
		uploadResp, err := session.UploadPart(&s3.UploadPartInput{
			Body:          bytes.NewReader(fileBytes),
			Bucket:        resp.Bucket,
			Key:           resp.Key,
			PartNumber:    aws.Int64(int64(partNum)),
			UploadId:      resp.UploadId,
			ContentLength: aws.Int64(int64(len(fileBytes))),
		})
		// Upload failed
		if err != nil {
			// fmt.Println(err)
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

func uploadObject(session *s3.S3, filename string) (resp *s3.PutObjectOutput) {
	bucket := GetEnvWithKey("AWS_S3_BUCKET")

	file, err := os.Open(filename)
	defer file.Close()
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()

	if err != nil {
		// fmt.Println(err)
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

	resp, _ = session.PutObject(params)

	return resp
}

func RemoveTempFile(fd *os.File) {
	for _, err := range [...]error{fd.Close(), os.Remove(fd.Name())} {
		if err != nil {
			log.Print("failed to remove temp file: ", err)
		}
	}
}
