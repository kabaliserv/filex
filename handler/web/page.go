package web

import (
	"github.com/kabaliserv/filex/web"
	"net/http"
)

func HandleIndex() http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {
		out := web.IndexFile()
		w.Write(out)
	}
	return f
}
