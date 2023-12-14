package main

import (
	"s3-queue-reader/handlers"
	"s3-queue-reader/infrastructures"
)

func main() {
	dbConnection := infrastructures.InitDBConnection(infrastructures.GetMySQLConfig())
	mySQL := infrastructures.NewMySQLHandler(dbConnection)

	defer dbConnection.Close()

	queueServiceConfig := infrastructures.GetServiceProviderConnectionConfig()
	resultedQueueName := infrastructures.GetResultedQueueName()
	queues := []string{resultedQueueName}

	queueService := infrastructures.NewQueueService(queueServiceConfig, queues)

	dbWriter := handlers.NewDBWriter(mySQL, queueService, resultedQueueName)

	dbWriter.Run()
}
