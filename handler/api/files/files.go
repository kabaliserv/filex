package files

import (
	"github.com/kabaliserv/filex/core"
	tusd "github.com/tus/tusd/pkg/handler"
	"log"
)

func PreUploadCreate(files core.FileStore) func(hook tusd.HookEvent) error {
	return func(hook tusd.HookEvent) error {
		//return files.AddInCache(hook.Upload.ID)
		return nil
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
