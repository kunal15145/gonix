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

func processListFiles(allFiles []string, continuationToken *string) ([]string, string, error) {
	var options *s3.ListObjectsV2Input
	if *continuationToken == "" {
		options = &s3.ListObjectsV2Input{
			Bucket: aws.String(s3BucketName),
		}
	} else {
		options = &s3.ListObjectsV2Input{
			Bucket:            aws.String(s3BucketName),
			ContinuationToken: continuationToken,
		}
	}
	resp, err := s3Session.ListObjectsV2(options)
	if err != nil {
		fmt.Println(err)
		return []string{}, "", err
	}
	for _, item := range resp.Contents {
		allFiles = append(allFiles, *item.Key)
	}
	return allFiles, *resp.NextContinuationToken, nil
}

/*
ListAllFiles -> returns files present in the current directory
*/
func ListAllFiles(currentPath string) ([]string, error) {

	var allFiles []string
	var continuationToken string

	allFiles, continuationToken, err := processListFiles(allFiles, &continuationToken)
	if err != nil {
		return []string{}, err
	}
	for continuationToken != "" {
		allFiles, continuationToken, err = processListFiles(allFiles, &continuationToken)
		if err != nil {
			return []string{}, err
		}
	}
	return allFiles, nil
}
