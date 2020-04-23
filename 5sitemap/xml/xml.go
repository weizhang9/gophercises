package xml

import (
	"log"
	"encoding/xml"

	"sitemap/sitemap"
)

const version string = "http://www.sitemaps.org/schemas/sitemap/0.9"
// Marshal func
func Marshal(s sitemap.Sitemap) []byte {
	var x XMLmap
	var u url
	for _, v := range s.Links {
		u.Loc = v
		x.URLs = append(x.URLs, u)
	}
	x.Xmlns = version
	res, err := xml.Marshal(x)
	if err != nil {
		log.Fatalln(err)
	}
	return res
}

// XMLmap struct
type XMLmap struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs     []url  `xml:"url"`
} 


type url struct {
	Loc  string `xml:"loc"`
}

// <?xml version="1.0" encoding="UTF-8"?>
// <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
//   <url>
//     <loc>http://www.example.com/</loc>
//   </url>
//   <url>
//     <loc>http://www.example.com/dogs</loc>
//   </url>
// </urlset>