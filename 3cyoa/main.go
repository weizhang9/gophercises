package main

import (
	"flag"
	"fmt"
	"net/http"
	"log"
	"os"

	"create-your-own-adventure/cyoa"
)

func main() {
	port := flag.Int("port", 3000, "the port to start the app")
	filename := flag.String("file", "gophers.json", "JSON file with the CYOA story")
	flag.Parse()

	f, err := os.Open(*filename)
	if err != nil {
		log.Fatalln(err)
	}

	story, err := cyoa.JSONStory(f)
	if err != nil {
		log.Fatalln(err)
	}

	h := cyoa.NewHandler(story)
	fmt.Println("starting server on port", *port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}