package main

import (
	"s3-queue-reader/handlers"
	"s3-queue-reader/infrastructures"
	"time"
)

func main() {
	var s3Watcher *handlers.S3Watcher

	dbConnection := infrastructures.InitDBConnection(infrastructures.GetMySQLConfig())
	mySQL := infrastructures.NewMySQLHandler(dbConnection)

	defer dbConnection.Close()

	s3Config := infrastructures.GetS3ConnectionConfig()
	s3 := infrastructures.NewS3Handler(s3Config)

	queueServiceConfig := infrastructures.GetServiceProviderConnectionConfig()
	queues := []string{"one", "two", "three"}
	queueService := infrastructures.NewQueueService(queueServiceConfig, queues)

	s3Watcher = handlers.NewS3Watcher(mySQL, s3, infrastructures.GetConfigBucketName(), queueService)

	for {
		s3Watcher.Run()
		time.Sleep(time.Second * 3)
	}

}
