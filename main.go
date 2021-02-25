package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kunal15145/gonix/awshelper"
)

var currentDirectoryContext string = "/"

func init() {

	var bucketName string
	var profileName string
	var regionName string

	flag.StringVar(&bucketName, "bucket", "", "A valid AWS bucket name")
	flag.StringVar(&profileName, "profile", "", "A valid AWS profile set in ~/.aws/credentials")
	flag.StringVar(&regionName, "region", "", "s3 bucket region name")
	flag.Parse()

	if bucketName == "" || profileName == "" || regionName == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	err := awshelper.InitializeSession(profileName, regionName, bucketName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleLsCommand() ([]string, error) {
	return awshelper.ListAllFiles(currentDirectoryContext)
}

func main() {

	var command string

	for {

		fmt.Print("gonix :> ")
		fmt.Scanln(&command)

		switch command {
		case "exit":
			os.Exit(0)
		case "ls":
			_, err := handleLsCommand()
			if err != nil {
				fmt.Println("Error Listing files in the current directory")
			}
		case "pwd":
			fmt.Println(currentDirectoryContext)
		default:
			continue
		}
	}
}
