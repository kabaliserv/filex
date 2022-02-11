package upload

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/service/token"
	log "github.com/sirupsen/logrus"
	tusd "github.com/tus/tusd/pkg/handler"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type uploadHandler struct {
	accesses core.AccessUploadStore
	handler  *tusd.UnroutedHandler
	options  core.UploadOption
	files    core.FileStore
	users    core.UserStore
	tokens   token.Manager
	uploads  core.UploadStore
	storages core.StorageStore
}

func NewUploadRouter(
	options core.UploadOption,
	files core.FileStore,
	users core.UserStore,
	tokens token.Manager,
	uploads core.UploadStore,
	accesses core.AccessUploadStore,
	storages core.StorageStore,
	fileStorage core.FileStoreComposer,
) *uploadHandler {

	u := uploadHandler{
		options:  options,
		files:    files,
		users:    users,
		tokens:   tokens,
		uploads:  uploads,
		accesses: accesses,
		storages: storages,
	}

	handler, err := tusd.NewUnroutedHandler(tusd.Config{
		BasePath:                  "/api/upload",
		StoreComposer:             fileStorage.StoreComposer,
		PreUploadCreateCallback:   u.preUploadCreate,
		PreFinishResponseCallback: u.preUploadFinish,
		NotifyCreatedUploads:      true,
	})

	if err != nil {
		log.Panic(err)
	}

	u.handler = handler

	return &u
}

func (u *uploadHandler) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(u.uploadCtx)

	r.Post("/", u.handler.PostFile)
	r.Patch("/{uploadId}", u.handler.PatchFile)
	r.Head("/{uploadId}", u.handler.HeadFile)

	u.notifyHandle()

	return r
}

func (u *uploadHandler) uploadCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Del("Filex-User-UUID")
		r.Header.Del("Filex-Storage-ID")
		r.Header.Del("Filex-Access-Id")

		user, ok := r.Context().Value(core.User{}).(core.User)
		if ok {
			r.Header.Set("Filex-User-UUID", user.UUID)
			r.Header.Set("Filex-Storage-ID", strconv.FormatUint(uint64(user.Storage.ID), 10))
		}

		next.ServeHTTP(w, r)
	})
}

func (u *uploadHandler) preUploadCreate(hook tusd.HookEvent) error {
	upload := hook.Upload
	req := hook.HTTPRequest

	delete(upload.MetaData, "storageId")

	authorization := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if authorization != "" {
		token, err := u.tokens.FromString(authorization)
		if err != nil || token.Subject() != "upload" {
			return tusd.NewHTTPError(errors.New("unauthorized"), http.StatusForbidden)
		}
		accessId, ok := token.PrivateClaims()[`accessId`].(string)
		if !ok || accessId == "" {
			return tusd.NewHTTPError(errors.New("unauthorized"), http.StatusForbidden)
		}

		access, err := u.accesses.FindById(accessId)
		if err != nil {
			return err
		} else if access == nil || access.IsRevoked {
			return tusd.NewHTTPError(errors.New("unauthorized"), http.StatusForbidden)
		}

		if upload.SizeIsDeferred {
			if !access.AllowDeferredSize {
				return tusd.NewHTTPError(errors.New("unauthorized"), http.StatusForbidden)
			}
		} else {
			if upload.Size > access.MaxSize {
				return tusd.NewHTTPError(errors.New("unauthorized"), http.StatusForbidden)
			}
		}

		if access.StorageId != "" {
			users, err := u.users.Find(core.User{Storage: core.Storage{UUID: access.StorageId}})
			if err != nil {
				return err
			} else if len(users) == 0 {
				return tusd.NewHTTPError(errors.New("unauthorized"), http.StatusForbidden)
			}

			user := users[0]

			if user.Storage.QuotaEnable {
				newSize := upload.Size + user.Storage.Size
				if newSize > user.Storage.Quota {
					return tusd.NewHTTPError(errors.New("unauthorized"), http.StatusForbidden)
				}
			}

			upload.MetaData["storageId"] = strconv.FormatUint(uint64(user.Storage.ID), 10)
		}

		access.IsRevoked = true

		if err := u.accesses.Save(access); err != nil {
			return err
		}

		req.Header.Set("Filex-Access-Id", accessId)

	} else if userId := req.Header.Get("Filex-User-UUID"); userId != "" {
		user, err := u.users.FindByUUID(userId)
		if err != nil {
			if errors.Is(err, core.ErrNotFound) {
				return tusd.NewHTTPError(errors.New("unauthorized"), http.StatusForbidden)
			}

			log.Errorf("preUploadCreate(find user in database): %s", err)
			return err
		}

		if user.Storage.QuotaEnable {
			newSize := upload.Size + user.Storage.Size
			if newSize > user.Storage.Quota {
				return tusd.NewHTTPError(errors.New("unauthorized"), http.StatusForbidden)
			}
		}

		req.Header.Set("Filex-Storage-ID", strconv.FormatUint(uint64(user.Storage.ID), 10))

	} else {
		return tusd.NewHTTPError(errors.New("unauthorized"), http.StatusForbidden)
	}

	return nil
}

func (u *uploadHandler) preUploadFinish(hook tusd.HookEvent) error {

	upload := hook.Upload

	files, err := u.files.Find(core.File{FileId: hook.Upload.ID})
	if err != nil {
		return err
	} else if len(files) == 0 {
		return tusd.NewHTTPError(errors.New("unauthorized"), http.StatusForbidden)
	}

	file := files[0]
	if file.IsFinal {
		return tusd.NewHTTPError(errors.New("unauthorized"), http.StatusForbidden)
	}

	storage, err := u.storages.FindByID(file.StorageID)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return tusd.NewHTTPError(errors.New("unauthorized"), http.StatusForbidden)
		}
		return err
	}

	file.FileId = upload.ID
	file.Name = upload.MetaData["filename"]
	file.Type = upload.MetaData["filetype"]
	file.Size = upload.Offset
	file.CreatedAt = time.Now()
	file.IsFinal = true

	storage.Size += file.Size

	if err = u.storages.Save(&storage); err != nil {
		return err
	}

	if err = u.files.Save(&file); err != nil {
		return err
	}

	return nil
}

func (u *uploadHandler) notifyHandle() {
	go func() {
		for {
			select {
			case hook := <-u.handler.CreatedUploads:
				upload := hook.Upload
				storageId, _ := strconv.ParseUint(hook.HTTPRequest.Header.Get("Filex-Storage-ID"), 10, 32)

				file := core.File{
					FileId:    hook.Upload.ID,
					Name:      upload.MetaData["filename"],
					Type:      upload.MetaData["filetype"],
					Size:      upload.Size,
					StorageID: uint(storageId),
				}

				if err := u.files.Create(&file); err != nil {
					log.Errorf("notifyHandle(create file cache): %s", err)
				}
			}
		}
	}()
}
