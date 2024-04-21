package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/abdulkaderm36/gophercises/url-short/main/db"
	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		redirect, ok := pathsToUrls[path]
		if !ok {
			fallback.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, redirect, http.StatusFound)
	}
}

func DBHandler(db db.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		redirect := []byte{}
		db.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("PathsToUrls"))
			redirect = b.Get([]byte(path))
			return nil
		})
		if redirect == nil {
			fallback.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, string(redirect), http.StatusFound)
	}
}

type PathUrl struct {
	Path string `yaml:"path" json:"path"`
	Url  string `yaml:"url"  json:"url"`
}

func parseYaml(yml []byte) ([]PathUrl, error) {
	pathUrls := []PathUrl{}
	err := yaml.Unmarshal(yml, &pathUrls)
	if err != nil {
		return nil, fmt.Errorf("Invalid YAML: %s\n", err.Error())
	}
	return pathUrls, nil
}

func parseJson(jsn []byte) ([]PathUrl, error) {
	pathUrls := []PathUrl{}
	if err := json.Unmarshal(jsn, &pathUrls); err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func buildMap(pathUrls []PathUrl) map[string]string {
	pathsToUrls := make(map[string]string)

	for _, v := range pathUrls {
		pathsToUrls[v.Path] = v.Url
	}

	return pathsToUrls
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseYaml(yml)
	pathsToUrls := buildMap(pathUrls)
	return MapHandler(pathsToUrls, fallback), err
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//	[
//	 path: /some-path
//	 url: https://www.some-url.com/demo
//	]
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseJson(jsn)
	pathsToUrls := buildMap(pathUrls)
	return MapHandler(pathsToUrls, fallback), err
}
