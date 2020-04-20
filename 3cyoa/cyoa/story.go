package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Choose Your Own Adventure</title>
</head>
<body>
	<section class="page">
		<h1>{{.Title}}</h1>
		{{range .Paragraphs}}
			<p>{{.}}</p>
		{{end}}
		<ul>
		{{range .Options}}
			<li>
				<a href="/{{.Chapter}}">{{.Text}}</a>
			</li>
		{{end}}
		</ul>
	</section>

	<style>
	body {
	  font-family: helvetica, arial;
	}
	h1 {
	  text-align:center;
	  position:relative;
	}
	.page {
	  width: 80%;
	  max-width: 500px;
	  margin: auto;
	  margin-top: 40px;
	  margin-bottom: 40px;
	  padding: 80px;
	  background: #FCF6FC;
	  border: 1px solid #eee;
	  box-shadow: 0 10px 6px -6px #797;
	}
	ul {
	  border-top: 1px dotted #ccc;
	  padding: 10px 0 0 0;
	  -webkit-padding-start: 0;
	}
	li {
	  padding-top: 10px;
	}
	a,
	a:visited {
	  text-decoration: underline;
	  color: #555;
	}
	a:active,
	a:hover {
	  color: #222;
	}
	p {
	  text-indent: 1em;
	}
  </style>
</body>
</html>
`

// Story mao
type Story map[string]Chapter

// Chapter struct
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// Option struct
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

// JSONStory func
func JSONStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

// Function Options - use functions to wrap params to modify/customise original (struct) type

// HandlerOptions type function options
type HandlerOptions func(h *handler)

// WithTemplate func allows custom override of default template
func WithTemplate(t *template.Template) HandlerOptions {
	return func(h *handler) {
		h.t = t
	}
}

// WithPathFn func allows custom override of default path design/look
func WithPathFn(fn func(r *http.Request) string) HandlerOptions {
	return func(h *handler) {
		h.pathFn = fn
	}
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	return path[1:]
}

// NewHandler func
func NewHandler(s Story, opts ...HandlerOptions) http.Handler {
	h := handler{s, tpl, defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s Story
	t *template.Template
	pathFn func(r *http.Request) string
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := defaultPathFn(r)
	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			log.Fatalln(err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}