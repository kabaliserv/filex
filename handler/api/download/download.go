package download

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/service/token"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

var (
	reExtractFileID  = regexp.MustCompile(`([^/]+)\/?$`)
	reForwardedHost  = regexp.MustCompile(`host=([^;]+)`)
	reForwardedProto = regexp.MustCompile(`proto=(https?)`)
	reMimeType       = regexp.MustCompile(`^[a-z]+\/[a-z0-9\-\+\.]+$`)
)

func HandleDownloadFile(files core.FileStore, tokens token.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uploadId := chi.URLParam(r, "uploadId")
		if uploadId == "" {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
		query := r.URL.Query()
		t := query.Get("token")
		if t == "" {
			url := "/d?u=" + uploadId
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		}
		token, err := tokens.FromString(t)
		if err != nil || token.Subject() != "download file" {
			url := "/d?u=" + uploadId
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		}

		fileId := token.PrivateClaims()["fileId"].(string)
		file, err := files.FindByUUID(fileId)
		if err != nil {
			if errors.Is(err, core.ErrNotFound) {
				http.NotFound(w, r)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("downloadFile(find file in database): %s", err)
		}

		src, err := files.GetReader(&file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("downloadFile(get reader): %s", err)
			return
		}

		w.Header().Set("Content-Length", strconv.FormatInt(file.Size, 10))

		contentType, contentDisposition := filterContentType(&file)
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Disposition", contentDisposition)

		// If no data has been uploaded yet, respond with an empty "204 No Content" status.
		if file.Size == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		io.Copy(w, src)

		// Try to close the reader if the io.Closer interface is implemented
		if closer, ok := src.(io.Closer); ok {
			closer.Close()
		}
	}
}

func GetByFileCtx(files core.FileStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, ok := r.Context().Value(core.File{}).(core.File)
		if !ok {
			http.NotFound(w, r)
			return
		}

		downloadFile(w, r, files, &file)
	}
}

func downloadFile(w http.ResponseWriter, r *http.Request, files core.FileStore, file *core.File) {
	src, err := files.GetReader(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("downloadFile(get reader): %s", err)
		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(file.Size, 10))

	contentType, contentDisposition := filterContentType(file)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", contentDisposition)

	// If no data has been uploaded yet, respond with an empty "204 No Content" status.
	if file.Size == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	io.Copy(w, src)

	// Try to close the reader if the io.Closer interface is implemented
	if closer, ok := src.(io.Closer); ok {
		closer.Close()
	}
}

func filterContentType(file *core.File) (contentType string, contentDisposition string) {
	filetype := file.Type

	if reMimeType.MatchString(filetype) {
		contentType = filetype
		contentDisposition = "attachment"
	} else {
		// If the filetype from the metadata is not well formed, we use a
		// default type and force the browser to download the content.
		contentType = "application/octet-stream"
		contentDisposition = "attachment"
	}

	// Add a filename to Content-Disposition if one is available in the metadata
	if filename := file.Name; filename != "" {
		contentDisposition += ";filename=" + strconv.Quote(filename)
	}

	return contentType, contentDisposition
}
