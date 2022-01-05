package web

import (
	"bytes"
	"embed"
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"strings"
)

var (
	//go:embed dist
	embedWEB embed.FS

	dist fs.FS

	// Mode dev or prod
	Mode string

	indexFile *bytes.Buffer
)

type manifestJSON struct {
	Source sources `json:"main.ts"`
}

type sources struct {
	indexJS  string   `json:"file"`
	indexCSS []string `json:"css"`
}

func init() {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	var err error

	dist, err = fs.Sub(embedWEB, "dist")
	check(err)

	var indexJS string
	var indexCSS string

	if Mode == "prod" {
		var manifest manifestJSON

		src, err := dist.Open("mi")
		if err != nil {
			log.Fatal(err)
		}
		defer src.Close()
		if err := json.NewDecoder(src).Decode(&manifest); err != nil {
			log.Fatal(err)
		}
		indexJS = manifest.Source.indexJS
		indexCSS = manifest.Source.indexCSS[0]
	}

	// Parse Index File
	data := struct {
		DevMode  bool
		IndexJS  string
		IndexCSS string
	}{
		DevMode:  Mode != "prod",
		IndexJS:  indexJS,
		IndexCSS: indexCSS,
	}

	tmpl, err := template.New("app").Parse(string(ReadFile("/index.tmpl")))
	check(err)

	indexFile = bytes.NewBuffer([]byte{})

	err = tmpl.Execute(indexFile, data)
	check(err)

}

type WebFileSystem struct {
	fs fs.FS
}

func ReadFile(name string) []byte {
	name = strings.TrimPrefix(name, "/")
	src, err := fs.ReadFile(dist, name)
	if err != nil {
		return []byte{}
	}

	return src
}

func IndexFile() []byte {
	return indexFile.Bytes()
}

func New() fs.FS {
	return fs.FS(WebFileSystem{dist})
}

func (s WebFileSystem) Open(name string) (fs.File, error) {
	name = strings.TrimPrefix(name, "/")
	return s.fs.Open(name)
}
