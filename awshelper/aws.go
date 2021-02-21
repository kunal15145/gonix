package awshelper

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3Session *s3.S3
var s3BucketName string

/*
InitializeSession -> starts an aws s3 session
*/
func InitializeSession(profileName string, awsRegion string, bucketName string) error {

	s3BucketName = bucketName
	awsSession, error := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewSharedCredentials("", profileName),
	})
	if error != nil {
		return error
	}
	s3Session = s3.New(awsSession)
	_, hbErr := s3Session.HeadBucket(&s3.HeadBucketInput{Bucket: &bucketName})
	if hbErr != nil {
		return hbErr
	}
	return nil
}

/*
ListFiles -> returns files present in the current directory
*/
func ListFiles(currentPath string) (int, bool) {
	resp, err := s3Session.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(s3BucketName)})
	if err != nil {
		fmt.Println(err)
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
