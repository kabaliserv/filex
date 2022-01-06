package core

type StoreOption struct {
	DatabaseDriver      string
	DatabaseEndpoint    string
	FileStoreS3Bucket   string
	FileStoreS3Endpoint string
	FileStoreLocalPath  string
	SessionSecret       string
}
