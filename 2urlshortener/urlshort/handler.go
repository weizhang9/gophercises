
package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if val, ok := pathsToUrls[r.URL.RequestURI()]; ok {
			http.Redirect(w, r, val, 302)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathsToUrls, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}
	maps := buildMaps(pathsToUrls)
	return MapHandler(maps, fallback), nil
}

type pathToUrl struct {
	Path string `yaml:"path"`
	URL string `yaml:"url"`
}

func parseYaml(yml []byte) ([]pathToUrl, error) {
	var s []pathToUrl
	err := yaml.Unmarshal(yml, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func buildMaps(pathsToUrls []pathToUrl) map[string]string {
	m := make(map[string]string)
	for _, v := range pathsToUrls {
		m[v.Path] = v.URL
	}
	return m
}