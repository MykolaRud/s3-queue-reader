package interfaces

import "s3-queue-reader/models"

type IS3Handler interface {
	ListFiles(bucket string) []models.S3RemoteFile
	DownloadFile(bucket, keyFile, dest string) bool
}
