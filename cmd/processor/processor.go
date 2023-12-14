package main

import (
	"s3-queue-reader/handlers"
	"s3-queue-reader/infrastructures"
)

func main() {
	queueServiceConfig := infrastructures.GetServiceProviderConnectionConfig()
	processorQueueName := infrastructures.GetProcessorQueueName()
	resultedQueueName := infrastructures.GetResultedQueueName()
	queues := []string{processorQueueName, resultedQueueName}

	queueService := infrastructures.NewQueueService(queueServiceConfig, queues)
	queueProcessor := handlers.NewQueueProcessor(queueService, processorQueueName, resultedQueueName)

	queueProcessor.Run()
}
