package awshelper

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type FileSystemObject struct {
	filename    string
	isFile      bool
	isDirectory bool
}

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

func listFiles(allFiles []string, continuationToken *string) ([]string, string, error) {
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
	if resp.NextContinuationToken == nil {
		return allFiles, "", nil
	}
	return allFiles, *resp.NextContinuationToken, nil
}

/*
ListAllFiles -> returns files present in the current directory
*/
func ListAllFiles(currentPath string) ([]FileSystemObject, error) {

	var allFiles []string
	var continuationToken string

	allFiles, continuationToken, err := listFiles(allFiles, &continuationToken)
	if err != nil {
		return []FileSystemObject{}, err
	}
	for continuationToken != "" {
		allFiles, continuationToken, err = listFiles(allFiles, &continuationToken)
		if err != nil {
			return []FileSystemObject{}, err
		}
	}
	return processAllFiles(allFiles, currentPath), nil
}

func processAllFiles(allFiles []string, currentWorkingDirectory string) []FileSystemObject {

	var fileSystemArrayObject []FileSystemObject

	for _, file := range allFiles {
		fileSplitted := strings.Split(string("/")+file, currentWorkingDirectory)
		fileRelativelocation := fileSplitted[1]
		fileRelativeLocationSplitted := strings.Split(fileRelativelocation, "/")
		if len(fileRelativeLocationSplitted) >= 2 {

			// It's a directory
			fileSystemArrayObject = append(fileSystemArrayObject, FileSystemObject{
				filename:    fileRelativeLocationSplitted[0],
				isFile:      false,
				isDirectory: true,
			})

		} else if len(fileRelativeLocationSplitted) == 1 {

			// It's a file
			fileSystemArrayObject = append(fileSystemArrayObject, FileSystemObject{
				filename:    fileRelativeLocationSplitted[0],
				isFile:      true,
				isDirectory: false,
			})

		} else {
			panic("Fatal Error while processing file system")
		}
	}
	fmt.Println(fileSystemArrayObject)
	return fileSystemArrayObject
}
