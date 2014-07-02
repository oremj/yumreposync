package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/oremj/yumreposync/client"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("Usage: yumreposync package1 [package2...]")
	}

	err := client.Push(*server+"/publish", flag.Args())
	if err != nil {
		log.Fatal(err)
	}
}
