package main

import (
	"errors"
	"flag"
	"log"
	"net/http"

	"github.com/oremj/yumreposync/server"
)

func validateFlags() error {
	if *repoDir == "" {
		return errors.New("-repodir is required.")
	}

	if *s3Bucket == "" {
		return errors.New("-bucket is required.")
	}
	return nil
}

func main() {
	flag.Parse()
	if err := validateFlags(); err != nil {
		flag.Usage()
		log.Fatal(err)
	}

	http.HandleFunc("/publish", server.Publish)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
