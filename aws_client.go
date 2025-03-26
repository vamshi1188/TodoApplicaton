package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	bucketName = "vamshigodev"
	objectKey  = "tasks.json"
	s3Client   *s3.Client
	uploader   *manager.Uploader
	downloader *manager.Downloader
)

func initS3Client() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	if err != nil {
		log.Fatalf("Unable to load AWS SDK config: %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
	uploader = manager.NewUploader(s3Client)
	downloader = manager.NewDownloader(s3Client)
}

func loadTasksFromS3() {
	buf := manager.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(context.TODO(), buf, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})
	if err != nil {
		log.Printf("Failed to download tasks from S3: %v", err)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &tasks); err != nil {
		log.Printf("Failed to unmarshal tasks: %v", err)
		return
	}

	// Update taskIDCounter to avoid ID conflicts
	for _, task := range tasks {
		if task.ID >= taskIDCounter {
			taskIDCounter = task.ID + 1
		}
	}
}

func saveTasksToS3() {
	data, err := json.Marshal(tasks)
	if err != nil {
		log.Printf("Failed to marshal tasks: %v", err)
		showMessage("Failed to save tasks.")
		return
	}

	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      &bucketName,
		Key:         &objectKey,
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		log.Printf("Failed to upload tasks to S3: %v", err)
		showMessage("Failed to save tasks.")
		return
	}

	showMessage("Tasks saved successfully!")
}
