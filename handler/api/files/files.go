package files

import (
	"errors"
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/service/token"
	log "github.com/sirupsen/logrus"
	tusd "github.com/tus/tusd/pkg/handler"

	"net/http"
)

func PreUploadCreate(tokenManager token.Manager, options core.UploadOption) func(hook tusd.HookEvent) error {
	return func(hook tusd.HookEvent) error {
		upload := hook.Upload
		req := hook.HTTPRequest

		token, err := tokenManager.GetTokenFromHeaders(req.Header)
		if err != nil {
			log.Error(err)
			return tusd.NewHTTPError(errors.New("internal error"), http.StatusInternalServerError)
		}

		if token != nil {
			if sub := token.Subject(); sub == "upload file" {
				return nil
			}
		}

		if !upload.SizeIsDeferred && options.GuestAllow && upload.Size <= options.GuestMaxUploadSize {
			return nil
		}

		return tusd.NewHTTPError(errors.New("unauthorized"), http.StatusUnauthorized)
	}
}

func PreUploadFinish(files core.FileStore) func(hook tusd.HookEvent) error {
	return func(hook tusd.HookEvent) error {
		log.Printf("Upload Files: %v", hook.Upload.ID)
		file, err := files.New(hook.Upload.ID)
		log.Println(file)
		return err
	}
}
