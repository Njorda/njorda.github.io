package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"golang.design/x/clipboard"
)

var (
	platform string
	bucket   string
	region   string
	filePath string
	ttl      int
)

func parseFlags() {
	flag.StringVar(&platform, "platform", "aws", "The platform to store the file")
	flag.StringVar(&filePath, "filePath", "", "The file path")
	flag.StringVar(&bucket, "bucket", "", "The bucket to use")
	flag.StringVar(&region, "region", "eu-west-1", "The region to use")
	flag.IntVar(&ttl, "ttl seconds", 60, "The time to live for the tmp link in seconds")
	flag.Parse()
}

// UploadFile reads from a file and puts the data into an object in a bucket.
func uploadFile(client *s3.Client, bucketName string, objectKey string, fileName string) error {
	expires := time.Now().Add(time.Duration(2*ttl) * time.Second)
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("couldn't open file %v: %v",
			fileName, err)
	}
	defer file.Close()
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:  aws.String(bucketName),
		Key:     aws.String(objectKey),
		Body:    file,
		Expires: &expires,
	})
	if err != nil {
		return fmt.Errorf("couldn't upload file %v to %v:%v. Here's why: %v",
			fileName, bucketName, objectKey, err)
	}
	return nil
}

func presignDocument(ctx context.Context, presignClient *s3.PresignClient, bucket string, objectPath string) (string, error) {
	params := s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &objectPath,
	}
	req, err := presignClient.PresignGetObject(
		ctx,
		&params,
		func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(ttl) * time.Second
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to presign: %s", objectPath)
	}
	return req.URL, nil
}

func main() {
	ctx := context.Background()
	parseFlags()

	var url string
	switch platform {
	case "aws":
		// Load the Shared AWS Configuration (~/.aws/config)
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatal(err)
		}
		cfg.Region = region

		// Create an Amazon S3 service client
		client := s3.NewFromConfig(cfg)
		path := fmt.Sprintf("/tmp/%s", filePath)
		if err := uploadFile(client, bucket, path, filePath); err != nil {
			panic(err)
		}
		presignClient := s3.NewPresignClient(client)
		url, err = presignDocument(ctx, presignClient, bucket, path)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Shareable link copied to clipboard, valid: %v seconds \n", ttl)
	case "gcp":
		panic("gcp is not implemented")
	}

	if err := clipboard.Init(); err != nil {
		panic(err)
	}
	clipboard.Write(clipboard.FmtText, []byte(url))
}
