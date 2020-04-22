package main

import (
	"strings"
	"fmt"
	"log"

	"urlscaper/link"
)

var exampleHtml = `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">A link to another page</a>
  <a href="/link-two">A link to second page</a>
</body>
</html>
`

func main() {
	r := strings.NewReader(exampleHtml)
	links, err := link.Parse(r)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", links)
}