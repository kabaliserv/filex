package files

import (
	"errors"
	"fmt"
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/service/token"
	"github.com/lestrrat-go/jwx/jwt"
	gonanoid "github.com/matoous/go-nanoid"
	log "github.com/sirupsen/logrus"
	tusd "github.com/tus/tusd/pkg/handler"
	"strings"
	"time"

	"net/http"
)

type FileHandler struct {
	handler *tusd.UnroutedHandler
	options core.UploadOption
	files   core.FileStore
	users   core.UserStore
	tokens  token.Manager
}

func NewFilesHandler(
	options core.UploadOption,
	files core.FileStore,
	users core.UserStore,
	tokens token.Manager,
) *FileHandler {

	fileHandler := FileHandler{
		options: options,
		files:   files,
		users:   users,
		tokens:  tokens,
	}

	handler, err := tusd.NewUnroutedHandler(tusd.Config{
		BasePath:                "/api/files",
		StoreComposer:           files.TusdStoreComposer(),
		PreUploadCreateCallback: fileHandler.preUploadCreate(),
		NotifyCreatedUploads:    true,
		NotifyCompleteUploads:   true,
	})

	if err != nil {
		log.Panic(err)
	}
	fileHandler.handler = handler

	go func() {
		for {
			select {
			case hook := <-handler.CreatedUploads:
				contextUploadId := hook.HTTPRequest.Header.Get("filex-context-upload-id")
				cacheFile := core.FileCache{
					FileId:          hook.Upload.ID,
					ContextUploadID: contextUploadId,
				}
				_ = files.CreateInCache(&cacheFile)
			case hook := <-handler.CompleteUploads:
				_ = files.CreateFromCache(hook.Upload.ID)
			}
		}
	}()

	return &fileHandler
}

func (f *FileHandler) SecureUploadMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		fmt.Println(strings.Contains(r.URL.Path, "request"))

		if !strings.Contains(r.URL.Path, "request") && (r.Method == "PATCH" || r.Method == "POST") {
			t, err := f.tokens.GetTokenFromHeaders(r.Header)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Error(err)
				return
			} else if t == nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			contextUploadId := t.Issuer()
			r.Header.Del("filex-context-upload-id")
			r.Header.Set("filex-context-upload-id", contextUploadId)

			// check if the upload was initiated by the customer
			if r.Method == "PATCH" {
				filesCache, err := f.files.FindInCache(core.FileCache{ContextUploadID: contextUploadId})
				if err != nil || len(filesCache) == 0 {
					log.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				fileCache := filesCache[0]

				if fileCache == nil || fileCache.ContextUploadID != contextUploadId {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func (f *FileHandler) PostFile(w http.ResponseWriter, r *http.Request) {
	f.handler.PostFile(w, r)
}

func (f *FileHandler) GetAllFile(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}

func (f *FileHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	f.handler.GetFile(w, r)
}

func (f *FileHandler) HeadFile(w http.ResponseWriter, r *http.Request) {
	f.handler.HeadFile(w, r)
}

func (f *FileHandler) PatchFile(w http.ResponseWriter, r *http.Request) {
	f.handler.PatchFile(w, r)
}

func (f *FileHandler) DelFile(w http.ResponseWriter, r *http.Request) {
	f.handler.DelFile(w, r)
}

func (f *FileHandler) preUploadCreate() func(hook tusd.HookEvent) error {
	return func(hook tusd.HookEvent) error {
		upload := hook.Upload
		req := hook.HTTPRequest

		t, err := f.tokens.GetTokenFromHeaders(req.Header)
		if err != nil || t == nil || t.Subject() != "upload file" {
			if err != nil {
				log.Error(err)
			}
			return tusd.NewHTTPError(errors.New("unauthorized"), http.StatusUnauthorized)
		}

		fmt.Printf("%#v", t.PrivateClaims())

		storageId, ok := t.PrivateClaims()["storage_id"].(string)
		if ok {
			if upload.MetaData["storage_id"] != storageId {
				return tusd.NewHTTPError(errors.New("unauthorized"), http.StatusUnauthorized)
			}
		}

		//#######################################
		// Verification des droits utilisateurs... A FAIRE !!
		//#######################################

		if !upload.SizeIsDeferred && f.options.GuestAllow && upload.Size <= f.options.GuestMaxUploadSize {
			return nil
		}

		return tusd.NewHTTPError(errors.New("unauthorized"), http.StatusUnauthorized)
	}
}

//func PreUploadFinish(files core.FileStore) func(hook tusd.HookEvent) error {
//	return func(hook tusd.HookEvent) error {
//		log.Printf("Upload Files: %v", hook.Upload.ID)
//		file, err := files.New(hook.Upload.ID)
//		log.Println(file)
//		return err
//	}
//}

func (f *FileHandler) HandlerRequestUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := f.tokens.NewToken()

	token.Set(jwt.SubjectKey, `upload file`)
	token.Set(jwt.ExpirationKey, time.Now().Add(time.Second*5))

	if user, ok := ctx.Value(core.User{}).(*core.User); ok && user != nil {
		token.PrivateClaims()["storage_id"] = user.Storage.UUID.String()
	}

	// generate contextUploadId
	contextUploadId := gonanoid.MustID(10)
	token.Set(jwt.IssuerKey, contextUploadId)

	payload, err := f.tokens.Sign(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
