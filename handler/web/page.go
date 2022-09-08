package web

import (
	dist "github.com/kabaliserv/filex/web"
	"net/http"
)

func HandleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//if  r.URL.Path != "/" {
		//	http.Redirect(w, r, "/", 303)
		//	return
		//}

		out := dist.IndexFile()
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.Write(out)
	}
}
