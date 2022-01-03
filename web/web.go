package web

import (
	"bytes"
	"embed"
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

func init() {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	var err error

	dist, err = fs.Sub(embedWEB, "dist")
	check(err)

	// Parse Index File
	data := struct {
		DevMode bool
	}{
		DevMode: Mode != "prod",
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
