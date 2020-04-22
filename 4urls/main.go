package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	// "strings"

	"urlscaper/link"
)

var exampleHtml = `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">A link to another page 
	<span>some span</span>
  </a>
  <a href="/link-two">A link to second page</a>
</body>
</html>
`

func main() {
	file, err := ioutil.ReadFile("ex4.html")
	if err != nil {
		log.Fatalln(err)
	}
	r := bytes.NewReader(file)
	// r := strings.NewReader(exampleHtml)
	links, err := link.Parse(r)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", links)
}