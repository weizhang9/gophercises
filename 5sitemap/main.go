package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"sitemap/sitemap"
	"sitemap/xml"
)

func main() {
	entry := flag.String("url", "", "main entrypoint for sitemap builder")
	flag.Parse()
	if !strings.HasPrefix(*entry, "http") {
		log.Fatalln("Please provide an entrypoint url with protocol scheme")
	}

	s := sitemap.Parse(*entry)
	res := xml.Marshal(s)

	fmt.Println(string(res))
}

