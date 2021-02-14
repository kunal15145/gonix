package main

import (
	"flag"
	"fmt"
	"os"
)

var bucketName string
var accessKeyID string
var secretAccessKey string

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

func main() {

}
