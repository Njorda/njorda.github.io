---
layout: post
title: "ShareIt, a cli for sharing files"
subtitle: "Build a small go cli"
date: 2024-04-17
author: "Niklas Hansson"
URL: "/2024/04/17/share_it"
---

TLDR: code can be found [here](https://github.com/Njorda/homebrew-tools).

In this blog post we will dive down in to how to build a small CLI for sharing files. The goal is to go over how to build a go cli for sharing files. We will set it up so that a shareable link will be created with a set expiration time and the object will be cleaned up after twice that time. For a file we will do the following: 

1) Upload the object to cloud blob storage with a expiration time so it will be deleted.
2) Create a presigned shareable link, any one with the link will be able to download the file
3) Copy the link to the clip board to make it easy to share for anyone. 

When we have the CLI we will set up so it can be installed through brew. 


# Build the CLI 

In order to keep it simple and avoid to many deps we will use the go standard package flags package for command line arguments. 

```go 
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
```

we will run this as the first thing in the `main` functions. The next step is to create the client to interact with blob storage. We will use AWS and S3 to keep it simple but this could be implemented for any blob storage. 

```go 
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
```

The upload function if a slightly modified version of the AWS example: 

```go 
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
```

We can see that the expiration time of the file on s3 is twice the time to live(ttl) of the sharable link. We do not intend to be able to recreate the link twice but instead the user have to upload the file twice. The presign function is also similar to the AWS examples: 

```go 
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
```

And thus the complete CLI looks like this: 

```go
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
```

To run it you can do the following: 

```bash
go run main.go --platform aws --filePath test.png --bucket test-bucket --region eu-west-
```

or compile it: 

```bash
go build -o shareit .
chmod +x ./shareit
./shareit --platform aws --filePath test.png --bucket oskar-test-2 --region eu-west-1
```


# Create homebrew package

Homebrew offers `taps` as an alternative to add third party repositories. By default `tap` assumes the repository is on `GitHub`. We will combine this together with [goreleaser](https://github.com/goreleaser/goreleaser) to build binaries for several platforms and push them back to github. 

To get up and runing we need to do two things:

1) add a `.goreleaser.yml` [file at the root of the repo](https://github.com/Njorda/homebrew-tools/blob/main/.goreleaser.yml).
2) add the [github action](https://github.com/Njorda/homebrew-tools/blob/main/.github/workflows/gorelease.yaml) that will push the shareit.rb file back to the repo and build the binaries.  


Lets go over the different sections of the `.goreleaser.yml` and what id does: 


```yaml
brews:
- name: shareit
  homepage: https://github.com/Njorda/homebrew-tools
  repository:
    owner: Njorda
    name: homebrew-tools
    branch: main
builds:
- main: ./src/shareit/
```

Most of the yaml file is boiler plat from goreleases examples. The key differences are the values for `homepage`, `owner`(the org or person that has the repo) and `name`(the name of the repo). I mixed up the owner and forgot that I created it in an organisation and not on my private account. Another key difference from the examples are the `builds:main` which points to where the main go file lives. 


The next step is to update the github action  for the build, `.github/workflows/gorelease.yaml`. This is also boilerplate from the [goreleaser examples](https://goreleaser.com/ci/actions/#workflow).

The repos needs to be public to work. 

That should be it and we should not be ready to create a release doing the following: 


```bash
git tag v0.0.x
git push origin v0.0.x
```

Check your github actions builds, mine can be found [here](https://github.com/Njorda/homebrew-tools/actions). When the actions are done and hopefully succeded. You can install it doing the following: 


```bash
brew tap Njorda/tools
brew install ShareIt
```

You should then be able to run: 

```
shareit --help
```

And we are done. Hope you found it helpful and now know how to build a small CLI in go and how to distribute it with Homebrew. 
