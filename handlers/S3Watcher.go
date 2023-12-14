package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"s3-queue-reader/infrastructures"
	"s3-queue-reader/interfaces"
	"s3-queue-reader/models"
	"s3-queue-reader/repositories"
)

func NewS3Watcher(dbHandler interfaces.IDbHandler, s3Handler interfaces.IS3Handler, bucketName string, queueService interfaces.IQueueService) *S3Watcher {
	dbFilesRepo := repositories.DBFilesRepository{dbHandler}
	tempFileName := infrastructures.GetTempDirectory() + "/" + "watcher_temp.parquet"

	return &S3Watcher{
		DBHandler:         dbHandler,
		DBFilesRepository: dbFilesRepo,
		S3Handler:         s3Handler,
		BucketName:        bucketName,
		TempFileName:      tempFileName,
		QueueService:      queueService,
	}
}

type S3Watcher struct {
	DBHandler         interfaces.IDbHandler
	DBFilesRepository repositories.DBFilesRepository
	S3Handler         interfaces.IS3Handler
	BucketName        string
	TempFileName      string
	QueueService      interfaces.IQueueService
}

func (h *S3Watcher) Run() {
	fmt.Println("sync")

	//read S3 list of files
	files := h.S3Handler.ListFiles(h.BucketName)

	//detect new files
	newFiles := h.FilterNewFiles(files)

	//download file
	for _, file := range newFiles {
		//download file
		if !h.S3Handler.DownloadFile(h.BucketName, file.Name, h.TempFileName) {
			fmt.Println("Error downloading file " + file.Name)

			continue
		}

		queueName := filepath.Dir(file.Name)

		//push file data to the queue
		if h.PushFileToQueue(h.TempFileName, queueName) {
			//mark as processed to DB
			h.DBFilesRepository.SetAsProcessed(file)
		}
	}
}

func (h *S3Watcher) PushFileToQueue(tempFileName string, queueName string) bool {

	fmt.Println("push file", tempFileName, "into queue", queueName)

	byteData, err := os.ReadFile(tempFileName)
	if err != nil {
		fmt.Println("Couldn't read file", tempFileName)
		return false
	}

	err = h.QueueService.PushDataToQueue(byteData, queueName)
	if err != nil {
		fmt.Println("Couldn't push file", tempFileName, "to queue", queueName)
		return false
	} else {
		return true
	}
}

func (h *S3Watcher) FilterNewFiles(files []models.S3RemoteFile) []models.S3RemoteFile {
	filteredFiles := make([]models.S3RemoteFile, 0)
	for _, file := range files {
		if !h.DBFilesRepository.FileExists(file) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles
}
