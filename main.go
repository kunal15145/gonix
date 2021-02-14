package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var bucketName string
var accessKeyID string
var secretAccessKey string
var currentDirectoryContext string = "/"

func init() {

	// check for bucket name
	flag.StringVar(&bucketName, "bucket", "", "A valid AWS bucket name")
	flag.Parse()

	if bucketName == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	_, err := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	if !err {
		fmt.Fprintf(os.Stderr, "error: %v\n", "AWS secret access key not found.")
		os.Exit(1)
	} else {
		secretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	}

	_, err = os.LookupEnv("AWS_ACCESS_KEY_ID")
	if !err {
		fmt.Fprintf(os.Stderr, "error: %v\n", "AWS access key id not found.")
		os.Exit(1)
	} else {
		accessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	}

}

func handleLsCommand(svc *s3.S3) (int, bool) {
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucketName)})
	if err != nil {
		return 1, false
	}
	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}
	return 1, true
}

func main() {

	var command string

	awsSession, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		fmt.Println("Error creating aws session")
		os.Exit(1)
	}

	s3Session := s3.New(awsSession)

	for {

		fmt.Print("gonix :> ")
		fmt.Scanln(&command)

		switch command {
		case "exit":
			os.Exit(0)
		case "ls":
			_, status := handleLsCommand(s3Session)
			if status {
				fmt.Println("Error Listing files in the current directory")
			}
		case "pwd":
			fmt.Println(currentDirectoryContext)
		default:
			continue
		}
	}
}
