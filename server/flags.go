package main

import "flag"

var addr = flag.String("addr", ":8080", "Bind Address")
var repoDir = flag.String("repodir", "", "Directory of Repo. (required)")
var s3Bucket = flag.String("bucket", "", "An s3 bucket, maybe contain a subdirectory. (required)")
