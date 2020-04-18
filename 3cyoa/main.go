package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var tpl *template.Template

var helper = template.FuncMap{
	"pascalcase": pascalCase,
}

func init() {
	tpl = template.Must(template.New("index.gohtml").Funcs(helper).ParseFiles("templates/index.gohtml"))
}

func main() {
	data := parseJSON("gopher.json")
	http.HandleFunc("/", homeHandler(data.Arc["intro"]))
	for url, value := range data.Arc {
		if url == "intro" {
			continue
		}
		http.HandleFunc(fmt.Sprintf("/%v", url), homeHandler(value))
	}
	log.Println("listening on port 8088")
	log.Fatal(http.ListenAndServe(":8088", nil))
}

func homeHandler(data arc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, data)
	}
}

type arcFull struct {
	Arc map[string]arc `json:"games"`
} 

type arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
} 

func parseJSON(filename string) arcFull {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	var arcs arcFull
	json.Unmarshal(data, &arcs)
	return arcs
}

func pascalCase(s string) string {
	ss := strings.Split(s, "-")
	return strings.Title(strings.Join(ss, " "))
}