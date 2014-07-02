package main

import (
	"flag"
	"log"
	"net/http"
)

func validateFlags() {
	if *repoDir == "" {
		flag.Usage()
		log.Fatal("-repodir is required.")
	}

	if *s3Bucket == "" {
		flag.Usage()
		log.Fatal("-repobucket is required.")
	}
}

func main() {
	flag.Parse()
	validateFlags()
	log.Fatal(http.ListenAndServe(*addr, nil))
}
