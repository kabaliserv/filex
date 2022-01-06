package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/memorylocker"
	"github.com/tus/tusd/pkg/s3store"
)

type S3 struct {
	Bucket   string
	EndPoint string
	store    s3store.S3Store
	locker   *memorylocker.MemoryLocker
	composer *tusd.StoreComposer
}

func New(
	Bucket string,
	EndPoint string,
) *S3 {
	store := &S3{
		Bucket:   Bucket,
		EndPoint: EndPoint,
	}
	return store.init()
}

func (s *S3) init() *S3 {
	s.locker = memorylocker.New()
	s3Config := aws.NewConfig()
	if s.EndPoint != "" {
		s3Config = s3Config.WithEndpoint(s.EndPoint).WithS3ForcePathStyle(true)
	}

	s.store = s3store.New(s.Bucket, s3.New(session.Must(session.NewSession()), s3Config))
	s.store.ObjectPrefix = ""
	s.store.PreferredPartSize = 50 * 1024 * 1024
	s.store.DisableContentHashes = false

	return s
}

func (s *S3) GetStoreComposer() *tusd.StoreComposer {
	if s.composer == nil {
		s.composer = tusd.NewStoreComposer()

		s.store.UseIn(s.composer)
		s.locker.UseIn(s.composer)
	}

	return s.composer
}