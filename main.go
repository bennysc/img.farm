package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
)

var key = os.Getenv("SPACES_KEY")
var secret = os.Getenv("SPACES_SECRET")

var s3Config = &aws.Config{
	Credentials: credentials.NewStaticCredentials(key, secret, ""),
	Endpoint:    aws.String("https://nyc3.digitaloceanspaces.com"),
	Region:      aws.String("us-east-1"),
}

var newSession, s3Err = session.NewSession(s3Config)
var s3Client = s3.New(newSession)

func main() {

	if s3Err != nil {
		fmt.Printf(s3Err.Error())
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/upload", Upload)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func Upload(w http.ResponseWriter, r *http.Request) {
	spaces, err := s3Client.ListBuckets(nil)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	for _, b := range spaces.Buckets {
		fmt.Fprintf(w, aws.StringValue(b.Name)+"\n")
	}
}
