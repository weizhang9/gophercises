package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"sitemap/sitemap"
)

func main() {
	entry := flag.String("url", "", "main entrypoint for sitemap builder")
	flag.Parse()
	if !strings.HasPrefix(*entry, "http") {
		log.Fatalln("Please provide an entrypoint url with protocol scheme")
	}

	s := sitemap.Parse(*entry)
	fmt.Printf("%#v",s)
}

